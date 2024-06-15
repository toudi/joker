# Available commands

For the time being there's no fancy command-line parsing like cobra but I plan to add one soon

## `joker up`

The main command. It parses the jokerfile, creates the graph of dependencies and instantiates your project.

Once joker is up an running you should see the output from your service's stdout and stderr on your terminal. You can press Ctrl-C and joker will then proceed to stop the services (send shutdown signal to them)

## `joker clear`

This command performs the `cleanup` step of the service. This means that the next time you run `joker up`, the `bootstrap` step will be performed for the affected services. If you don't pass any arguments, joker will perform `cleanup` for all of possible services. You can also limit this to a service by passing the argument:

```shell
# all services
joker clear
# single service
joker clear my-service
```
