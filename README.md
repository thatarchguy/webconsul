# Webconsul

Web API to get a prefix from consul and serve


## Quick Example

This short example assumes Consul is installed locally.

1. Start a Consul cluster in dev mode.

    ```
    $ consul agent -dev
    ```

1. Write some data.
    ```
        $ consul kv put my-app/address 1.2.3.4
        $ consul kv put my-app/port 80
        $ consul kv put my-app/max_conns 5
    ```

1. build and run the server
    ```
    $ go build

    $ ./webconsul -addr http://localhost:8500
    ```

1. Get the keys
    ```
    $ curl localhost:8080/v1/consul?prefix=my-app
    address: 1.2.3.4
    max_conns: 5
    port: 80
    ```
1. Get a key
    ```
    $ curl localhost:8080/v1/consul?key=my-app/address
    address: 1.2.3.4

    $ curl localhost:8080/v1/consul?key=my-app/address&upper=1
    ADDRESS: 1.2.3.4

    $ curl 'localhost:8080/v1/consul?key=my-app/address&override=bind-address'
    bind-address: 1.2.3.4
    ```
