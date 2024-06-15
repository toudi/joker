# What this project is **NOT** about

If it wasn't clear from the start, let me clarify this. This project is **NOT** about installing the components that are being used inside your projects. In other words - if you know that you require to have, say, PostgreSQL or redis or python - it is your responsibility to install the packages on your OS. Theoretically, I could imagine some descriptive way of specifying the installation method and the required packages, something like:

```yaml
required-packages:
  items:
    - postgresql
    - redis
  package-manager: apt # yum / zypper / dnf / homebrew
```

but I don't like this. I think it could create more problems than it would try to solve. Like.. there are certain differences when it comes to installing, say, postgres on homebrew vs doing it on debian or .. slackware. Though I'm open to suggestions :-)
