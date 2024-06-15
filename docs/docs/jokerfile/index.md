# Jokerfile

Jokerfile is the heart of joker. It's a fairly simple yaml file and it defines all the steps that are required to instantiate your application, as well as it allows to define some common, shared values that you can then interpolate. Let's go over the sections of the example jokerfile.

## env

Example:

```yaml
env:
  data_dir: /home/toudi/joker/services/
  redis_port: 5555
  postgres:
    username: django-app
    password: test
    dbname: example-db
    port: 2222
```

This section defines a set of variables that can be interpolated against items in the `services` section. In this example, you can imagine that if you set up redis to be running at port `5555` then you can reference the said port via <code v-html="'{{ env.redis_port }}'"></code> so that you wouldn't have to copy and paste the value as it's error-prone.

::: info
What I plan to implement in the future is to move the `data_dir` to a separate property. this would allow joker to archive all the data files to a tar.gz file so that you could achieve portability.
:::

## services

An array of services that define your project. Please refer to [Service](/docs/jokerfile/service/index) documentation for in-depth documentation.
