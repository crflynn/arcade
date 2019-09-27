# golyglot

Golyglot is a hosting service for publishing static software documentation (sphinx, exdoc, cratedocs) from multiple languages.

To launch the service locally, clone the repo and run

```bash
docker build -t golyglot .
docker run -v $(pwd)/docs:/tmp/docs -p 6060:6060 -it golyglot
```

To see the home page visit

```bash
http://localhost:6060/
```

You will be prompted for a username and password. These are specified by environment variables along with the docs root directory and port.
To override the default environment variables use:

```bash
docker run -p 6060:6060 \
    -e GOLYGLOT_PORT='6060'
    -e GOLYGLOT_DOCROOT='/tmp/docs'
    -e GOLYGLOT_USERNAME='admin'
    -e GOLYGLOT_PASSWORD='password'
    -it golyglot
```

To push new documentation, tar the docs contents and submit a PUT request:

```bash
# Navigate into the docs build directory (where index.html is).
# This will differ depending on language
# but the point is that we want the static docs built.
# The example here uses the sphinx build dir for html
cd docs/_build/html
# create a tar file of the contents excluding the top folder
tar -czf ../myproject.tar.gz .
cd ..
curl -u username:password -X PUT 'http://localhost:6060/docs/myproject/latest' --upload-file myproject.tar.gz
curl -u username:password -X PUT 'http://localhost:6060/docs/myproject/v1.2.3' --upload-file myproject.tar.gz
```

To delete documentation

```bash
# delete a single version
curl -u username:password -X DELETE 'http://localhost:6060/docs/myproject/v1.2.3'
# delete an entire project
curl -u username:password -X DELETE 'http://localhost:6060/docs/myproject'
```
