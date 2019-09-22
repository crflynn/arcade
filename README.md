# golyglot

Golyglot is a static hosting service for publishing software documentation (sphinx, exdoc, cratedocs) from any language.

To see the home page

```bash
curl http://localhost:6060/
```

You will be prompted for a username and password on clicking the link to view documentation. These are specified by environment variables along with the docs root directory and port.

To push new documentation, tar the docs contents and submit a PUT request to

```bash
# go into the docs build directory (where index.html) is
# this will differ depending on language
# but the point is that we want the static docs built
cd docs/_build
# create a tar file of the contents excluding the top folder
tar -czf ../myproject.tar.gz .
cd ..
curl -u username:password -X PUT 'http://localhost:6060/myproject' --upload-file myproject.tar.gz
```

To delete documentation

```bash
curl -u username:password -X DELETE 'http://localhost:6060/myproject'
```

This will delete the project and its contents



