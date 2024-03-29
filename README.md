# watchclock

**NOTE: This software is not yet fully implemented, do not use it yet.** I am publishing it for documentation purposes.

`watchclock` renews S3 Object Locks which are about to expire.

Object Lock prevents attackers from deleting your S3 objects even if they've managed to steal AWS credentials. For each S3 object, you choose a date in the future, and the object cannot be deleted before that date.

But after the expiration date, your object becomes deletable again until you set a _new_ date.

To prevent this, [configure your AWS credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html)) and then run `watchclock` from the command line:

```sh
watchclock --minimum-days 8 --renew-for 30 my-bucket-name
```

For any object in `my-bucket-name` that has an Object Lock retention date less than 8 days from now, this will reset that date to 30 days from now.

By running this command once per week, you ensure that every object in the bucket will have Object Lock set to between 1 and 30 days in the future.

## Deleting Objects

To intentionally delete an Object-Locked item, you must let the lock expire first. Because Object Lock works only on versioned buckets, you can use a workflow like this:

- Delete the object, which will write a delete marker to the bucket. This will not delete the old version, because old versions are always preserved by default, and also because the old version has Object Lock enabled.

- Wait for Object Lock to expire on the old version.

- Delete the old version.

By default `watchclock` does not update locks on older versions of objects, so if an object is deleted its older versions will automatically age out of Object Lock, at which point they can be pruned by a lifecycle rule like this one:

```sh
aws s3api put-bucket-lifecycle-configuration --lifecycle-configuration file://lifecycle-rules.json --bucket mybucket
```

where the contents of lifecycle-rules.json is:

```json
{
  "Rules": [
    {
      "ID": "prune-deleted-objects",
      "Prefix": "",
      "Status": "Enabled",
      "NoncurrentVersionExpiration": {
        "NoncurrentDays": 30
      }
    }
  ]
}
```

## Object Versions

In a versioned bucket every _version_ of an object has its own Object Lock. By default `watchclock renew` only updates the lock on the _current_ version. To update Object Locks on _all_ versions of each object, pass the `--all-versions` option. 

The `--all-versions` option is _far_ less efficient to run, because it requires one API call per affected version.

`--all-versions` will not update Object Locks on currently-deleted objects -- that is, objects whose current version is a delete marker. To change this behavior, add an additional option, `--include-deleted-objects` -- this will renew Object Locks on older versions of deleted objects.

## Caching

Retrieving the Object Lock information for an S3 object requires a `GetObjectRetention` API call per object version, which is very slow. So if there is an existing DynamoDB table called `watchclock-cache` which `watchclock` can read and write, it will automatically cache Object Lock information in that table.

- Pass the `--use-cache` option to `watchclock` to have it create the DynamoDB table for you. Once the table exists, you don't need to provide the `--use-cache` option when running `watchclock` in the future -- it will use the cache automatically.

- Use the `--cache-table-name` option to change the name of the cache table in DynamoDB.

- To clear and rebuild the cache, pass the `--clear-cache` option.

- To disable the use of the cache, pass the `--no-cache` option.


### Use Caching Carefully

The `watchclock` cache is designed with the assumption that `watchclock` is the only tool which alters Object Locks after they are initially created. If some other tool modifies an Object Lock, the `watchclock` cache could produce incorrect behavior. For example:

- Suppose we store a new object named `thing` with a retention date 10 days in the future.
- We run `watchclock --minimum-days 4`, which notes that `thing` will need to be renewed 6 days from now and records this in the cache.
- If we write a new version of `thing` with a retention date only 2 days in the future, there will now be a window (between three and six days from now) when the Object Lock on `thing` has expired but `watchclock --minimum-days 4` does not renew it because its cached information is stale.
