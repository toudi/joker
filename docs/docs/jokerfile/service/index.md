# Service

It is an abstract concept that defines a component of your project. It can either be a program that is provided for you (like a database) or a program that you are developing (say - a python app or a go binary)

## Service dependencies

Let's consider the following example:

```yaml
services:
  redis:
    # ...
  django:
    depends:
      - redis
      - postgres
  postgres:
    # ...
```

I think that's pretty self-explainatory. joker parses the jokerfile and scans trough services. Then, it orders them by their position in the dependency graph (the `depends` keyword is optional). In the above example, it would instantiate services in the following order:

1. redis
2. postgres
3. django

The position of the services in the yaml file is irrelevant

## Execution context

Each of the following properties: `bootstrap`, `command`, `cleanup` can be defined in 3 possible ways:

```yaml
services:
  my-service:
    # as a string
    bootstrap: touch /tmp/some-file
  my-service-2:
    # as a list of strings
    bootstrap:
      - touch /tmp/some-file
      - rm -f /tmp/some-file
  my-service-3:
    # as a dictionary (useful when you want to override the environment variables)
    bootstrap:
      env:
        MY_ENV_VARIABLE: foo
      commands:
        - touch /tmp/some-file
        - rm -f /tmp/some-file
```

You can also interpolate values from the `env` section in any member section like so:

```yaml
env:
  my_var: my_value
services:
  my-service:
    bootstrap:
      env:
        MY_ENV_VAR: "{{ env.my_var }}"
      commands:
        - "touch /tmp/{{ env.my_var }}"
```

## Service bootstrapping

This property defines steps that are required to be performed **prior** to launching a service for the first time. For instance, when you're developing a python app it's customary that you create a virtualenv, then install the dependencies with pip. For a go project you'd probably do `go mod tidy`. I can also imagine dumping some environmental variables to files and so on. Here's how to do it:

```yaml
services:
  postgres:
    bootstrap:
      - "{{ env.postgres.bin_dir }}/pg_ctl initdb -D {{ env.data_dir }}/postgres"
      - sed -in 's/#port = 5432/port = {{ env.postgres.port }}/' {{ env.data_dir }}/postgres/postgresql.conf
```

## Service cleanup

Cleanup is the reverse process of bootstrapping. It is performed when you want to clear all the prepared data for a service. Again, for a django app this would mean removing the virtualenv, maybe removing the SQLite file (if you're using SQLite that is). Here are some examples:

```yaml
services:
  django-app:
    dir: "{{ env.data_dir }}/django-app" // [!code highlight]
    cleanup:
      - rm -rf .env envs
  postgres:
    cleanup:
      - rm -rf {{ env.data_dir }}/postgres
```

::: info
Because of the `dir` property in `django-app` service, you can specify the relative paths to the directories
:::

## Running the service

This is probably the only required property as it defines how your service should be instantiated.

::: warning
If your service was not bootstrapped then the `bootstrap` section will be interpreted and executed **before** step(s) in `command` section.
:::

```yaml
services:
  redis:
    dir: "{{ env.data_dir }}/redis"
    command: redis-server --port {{ env.redis_port }}
  django-app:
    dir: "{{ env.data_dir }}/django-app"
    command:
      - source .env/bin/activate
      - envdir envs python manage.py migrate
      - envdir envs python manage.py runserver
  postgres:
    command: /opt/homebrew/bin/postgres -D {{ env.data_dir }}/postgres
```

::: warning
Joker checks for the service liveness during the instantiation. This means that if any problem occurs during `joker up`, joker will stop all the running services and exit. I understand that this may cause a problem as maybe you'd want to define a service that is allowed to exit (and therefore opt-out from the liveness check). I'll try to implement it asap. for the time being, your best option is to add this to the `bootstrap` or `command` sections as a step.
:::

## Hot reloading

::: info
This functionality would not be possible without `fsnotify` package. A big thanks!
:::

If your framework can listen on file modifications by itself (like Django does for example) then you don't have to do anything. If, however you'd like the service to be automatically rebuilt when files are modified, use the following snippet:

```yaml
services:
  my-service:
    dir: /some/directory/to/watch
    hot-reload: true
```

Yes I understand that it looks counter-intuitive, however as I've stated before, sometimes you may want to define a `dir` property for the service to use relative paths.

Long story short, what this will do is it will add `/some/directory/to/watch` (and any of it's subdirectories) to notification watcher and when you'll save your source files, the service will be automatically rebuilt.

You can also use one of the alternative syntaxes:

```yaml
services:
  my-service:
    hot-reload: /some/directory/to/watch
  my-other-service:
    hot-reload:
      - /some/directory
      - /some/other/directory
```

Again, any subdirectories discovered upon inspection will be added to the watchlist.
