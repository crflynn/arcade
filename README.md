# arcade

[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/crflynn/arcade)](https://hub.docker.com/r/crflynn/arcade)

`arcade` is a (very) basic web service written in Go for publishing static software documentation (sphinx, exdoc, cratedocs) from multiple projects with multiple versions across multiple languages. The service is meant to be used by private teams and is (optionally) secured by basic auth.

To launch the service locally, clone the repo and run

```bash
make build
make serve
```

or using the Dockerhub image

```
docker pull crflynn/arcade
docker run -p 6060:6060 -it crflynn/arcade
```

To see the home page visit

```bash
http://localhost:6060/
```

If you are running from a cloned repo using ``make``, you will be prompted for a username and password (default is admin:admin from docker-compose.yml). These are specified by environment variables along with the docs root directory and port.
To override the default environment variables use:

```bash
docker run -p 6060:6060 \
    -e ARCADE_PORT='6060' \
    -e ARCADE_DOCROOT='/docs' \
    -e ARCADE_USERNAME='username' \
    -e ARCADE_PASSWORD='password' \
    -it crflynn/arcade
```

To push new documentation, tar the docs contents and submit a PUT request:

```bash
# Navigate into the docs build directory (where index.html is).
# This will differ depending on language
# but the point is that we assume the docs are already built.
# The example here uses the default sphinx build dir for html
cd docs/_build/html
# create a tar file of the contents excluding the top folder
tar -czf ../myproject.tar.gz .
cd ..
curl -u username:password -X PUT 'http://localhost:6060/docs/myproject/latest' --upload-file myproject.tar.gz
curl -u username:password -X PUT 'http://localhost:6060/docs/myproject/v1.2.3' --upload-file myproject.tar.gz
```

To delete documentation, submit a DELETE request to the project and/or version path which should be removed.

```bash
# delete a single version
curl -u username:password -X DELETE 'http://localhost:6060/docs/myproject/v1.2.3'
# delete an entire project
curl -u username:password -X DELETE 'http://localhost:6060/docs/myproject'
```
