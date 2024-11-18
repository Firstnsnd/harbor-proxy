# harbor-proxy
A demo example for setting up a proxy for the Docker registry Harbor.

# How to Run
1. Update Domain
    Replace the following line in your code:
    ```go
    var domain = "your_harbor_domain";
    ```
    with your actual domain.
2. Run the Application
   Execute the following command in your terminal:
    ```sh
    go run main.go
    ```
3. Edit `/etc/hosts`
   Now,you need to edit the `/etc/hosts` file using `vim` and add the following entry:
    ```
    127.0.0.1       dev.test.com
    ```
4. Update Docker Daemon Configuration
   Edit the `/etc/docker/daemon.json` file with vim and add the following under the insecure-registries option:
   ```json
   {
    "ipv6": false,
    "insecure-registries" : ["dev.test.com:8099"]
    }
   ```
   **Note:** Ensure that there is a comma after "ipv6": false.
5. Restart Docker
   Don’t forget to restart the Docker service:
   ```sh
    systemctl restart docker
   ```
6. Tag Your Image
   Use the following command to tag your image name format as `dev.test.com:8099/****:***`. Replace the parts as necessary:
   ```sh
   docker tag d2c94e258dcb dev.test.com:8099/g299_remote/ss_g612/redis:latest
   ```
7. Push and Pull Images
   Now you can push and pull images through your proxy service to Harbor. Use the following commands:
   Push Command：
    ```sh
    $ docker push dev.test.com:8099/g299_remote/ss_g612/redis:latest
    ```
    Expected Output:
    ```shell
    The push refers to repository [dev.test.com:8099/g299_remote/ss_g612/redis]
    950a085c0a1c: Layer already exists 
    5f70bf18a086: Layer already exists 
    e4dbf0bd9d9d: Layer already exists 
    15ef09f03230: Layer already exists 
    40710ab1222c: Layer already exists 
    a64e92ee1239: Pushed 
    9a978e3d8066: Pushed 
    8e2ab394fabf: Pushed 
    latest: digest: sha256:9e0dedfd09654001f77a0d1dbe981717a746ef4d253f5d119c50be28e100337f size: 1986
    ```
    Pull Command:
    ```shell
    $ docker pull dev.test.com:8099/g299_remote/ss_g612/mongo 
    ```
    Expected Output:
    
    ```shell
    Using default tag: latest
    latest: Pulling from g299_remote/ss_g612/mongo
    32b824d45c61: Pull complete 
    2ffb4886d703: Pull complete 
    f36248f814ef: Pull complete 
    0b39c673fa12: Pull complete 
    a6208dee2c0e: Pull complete 
    876eec4ef49d: Pull complete 
    07ed7efc9402: Pull complete 
    e5e3b551bf11: Pull complete 
    Digest: sha256:04c3087ed64ee0a408f09f26c2ef8e6c704b2e2482753661004214c7306ccc68
    Status: Downloaded newer image for dev.test.com:8099/g299_remote/ss_g612/mongo:latest
    dev.test.com:8099/g299_remote/ss_g612/mongo:latest
    ```

