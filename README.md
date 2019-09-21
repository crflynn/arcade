# golyglot

Golyglot is a static hosting service for publishing software documentation (sphinx, exdoc, cratedocs) from any language.

To view the list of projects available

```bash
curl http://localhost:6060/
```


To push new documentation, tar the docs contents and submit a PUT request to

```bash
# go into the docs build directory (where index.html) is
# this will differ depending on language
# but the point is that we want the static docs built
cd docs/_build
# create a tar file of the contents excluding the top folder
tar -czf ../myproject.tar.gz .
cd ..
curl -X PUT 'http://localhost:6060/myproject' --upload-file myproject.tar.gz
```

To delete documentation

```bash
curl -X DELETE http://localhost:6060/myproject
```

This will delete the project and its contents



