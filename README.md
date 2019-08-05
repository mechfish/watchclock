# watchclock

**NOTE: This software is not yet fully implemented, do not use it yet.** I am publishing it for documentation purposes.

`watchclock` renews S3 Object Locks which are about to expire.

Object Lock prevents attackers from deleting your S3 objects even if they've managed to steal AWS credentials. For each S3 object, you choose a date in the future, and the object cannot be deleted before that date. 

But after the expiration date, your object becomes deletable again until you set a _new_ date.

To prevent this, [configure your AWS credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html)) and then run `watchclock` from the command line:

```sh
watchclock renew --bucket my-bucket-name --minimum-days 8 --renew-for 30
```

For any object in `my-bucket-name` that has an Object Lock retention date less than 8 days from now, this command will reset that date to 30 days from now.

By running this command once per week, you ensure that no Object Lock in the bucket will ever expire, but also that the retention dates will all be less than 30 days in the future.

## Deleting Objects

Intentionally deleting an Object-Locked item requires letting the lock expire first. `watchclock` has a `delete` command to help with this:

```sh
watchclock delete my-bucket-name/path/to/object
```

This command schedules the specified S3 object for future deletion. It will actually be deleted once its Object Lock has expired.

To change your mind and un-schedule an object for deletion, run:

```sh
watchclock undelete my-bucket-name/path/to/object
```

(`watchclock` stores its list of items to delete in a `.watchclock-to-delete` directory at the top level of the S3 bucket. `watchclock` will not set or renew Object Locks on this directory or its contents, and it will periodically clean up the directory.)

## Object Versions

So far we've discussed Object Locks as if the apply to whole S3 objects, but in a versioned bucket every _version_ of an object has its own Object Lock.

By default `watchclock renew` updates the lock on the _current_ version of an object but not on older versions. But by passing the `--all-versions` option, it will update the lock for _all_ versions of the object. Note that this is _much_ less efficient to run, because it requires one API call per affected version.

## Caching

Retrieving the Object Lock information for an S3 object requires one `GetObjectRetention` API call per object version, which is very slow. So, by default, `watchclock` caches Object Lock information in a DynamoDB table.

- `watchclock` will automatically create and use a DynamoDB table called `watchclock-cache`; use the `--cache-table` option to change this name.

- To clear and rebuild the cache, pass the `--clear-cache` option.

- To disable the use of the cache, pass the `--no-cache` option.

Caching is designed with the assumption that `watchclock` is the only tool which alters the locks after they are initially created. If some other tool modifies an Object Lock, the `watchclock` cache could produce incorrect behavior. For example:

- Suppose we store a new object named `thing` with a retention date 10 days in the future.
- We run `watchclock --minimum-days 4`, which notes that `thing` will need to be renewed 6 days from now and caches this fact.
- We write a new version of `thing` with a retention date only 2 days in the future. There will now be a window (between three and six days from now) when the Object Lock on `thing` has expired but `watchclock --minimum-days 4` does not renew it because its cached information is stale.
