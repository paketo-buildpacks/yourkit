# `paketobuildpacks/yourkit`
The Paketo Buildpack for YourKit is a Cloud Native Buildpack that contributes the [YourKit][y] Agent and configures it to
connect to the service.

[y]: https://www.yourkit.com

## Behavior

This buildpack will participate if all the following conditions are met

* `$BP_YOURKIT_ENABLED` is set

The buildpack will do the following at build time:

* Contributes a Java agent to a layer
* Contributes a helper that configures the agent at runtime

The helper binary runs at launch time and executes before your application. It reads the `BPL_*` configuration settings and uses them to configure the YourKit agent.

## Configuration

| Environment Variable | Description
| -------------------- | -----------
| `$BP_YOURKIT_ENABLED` | Whether to contribute YourKit support
| `$BPL_YOURKIT_ENABLED` | Whether to enable YourKit support
| `$BPL_YOURKIT_PORT` | Configure the port the YourKit agent will listen on. Defaults to `10001`.
| `$BPL_YOURKIT_SESSION_NAME` | Configure the session's name.

## Bindings

The buildpack optionally accepts the following bindings:

### Type: `dependency-mapping`

|Key                   | Value   | Description
|----------------------|---------|------------
|`<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>`

## Publishing the Port

When starting an application with the YourKit Profiler enabled, a port must be published.  To publish the port in Docker, use the following command:

```bash
$ docker run --publish <LOCAL_PORT>:<REMOTE_PORT> ...
```

The `REMOTE_PORT` should match the `port` configuration for the application (`10001` by default).  The `LOCAL_PORT` can be any open port on your computer, but typically matches the `REMOTE_PORT` where possible.

Once the port has been published, your YourKit Profiler should connect to `localhost:<LOCAL_PORT>` for profiling.

![YourKit Configuration](yourkit.png)

## License

This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
