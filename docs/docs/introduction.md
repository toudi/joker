# Why even bother when there's docker?

::: info
Yes, I know about podman. This is just a linguistic shortcut so that I wouldn't have to type docker/podman/oci-compliant-engine each time.
:::

I thought that I'd start the description of the project by addressing the elephant in the room: what this project is all about? Like.. what is the reason to use it when there's a perfectly valid solution out there and it's called docker. The way I see it - two things - simplicity and portability.

Let me reiterate: I've got absolutely nothing against docker. I think it's outstanding for production purposes - when you want to sandbox your application, when you want to containerize it, or when the dependencies are complicated to install to make the application work correctly (for instance, I could imagine using it for processing LaTeX documents). These are all perfectly valid arguments for using this tool. I just do not believe they are valid for local development.

This whole thing started from a chat that I had with my friends where we were discussing our workflow with regards to local development. Usually when I code, I like to open up a terminal and I usually have a set of commands that I issue in order to bring the local development stack up. For instance, for a python project I'd usually do:

```shell
source .env/bin/activate
envdir envs python ...
```

For redis it's as simple as

```shell
redis-server
```

For a go project you'd have to do

```shell
go build
./some-binary
```

Then for a database I'd instanciate a local server

That's good and all but as the project gets more complex you have to have a system in place to share the config values, maybe set up some local accounts and so on. Again, you can do all of this but it's tedious.

::: warning
The above commands describe an optimistic scenario in which you've already bootstrapped the project. Meaning - again - if we stick just to python example - that you've run a procedure similar to:

```shell
python3 -m venv .env
source .env/bin/activate
pip install -r requirements.txt
```

So now if you'd want to keep it simple you'd probably have to maintain a bash script that would check if the requirements were installed and if not, it would install them, otherwise it'd just run the project.
:::

To top it all off I'd want to archive the whole thing. Meaning - given a root path for the data files I'd want to tar.gzip it to be able to restore it on another computer. For the database this means that what I want is the copy of the database binary files - _NOT_ the SQL dump of the data. Why waste the time on exporting and importing the data I already have on my disk?

Then my friends said: well, why don't you just use docker? Well, because I don't want to, for couple of reasons. First of all, I would usually have the required components installed anyway (like postgres or redis). Second of all, docker has to sync the files onto the running image. But the files are already on my computer. This has a deeper consequence. If I'm debugging, I want the state of the application to reflect the code that I see in my editor. I don't want to rebuild the image just because I'm debugging. Also, on MacOS docker runs through a virtual machine which penalizes the performance. That does not even address the issue of having to monitor the disk space and removing old images to reclaim free space.

> Things should be made as simple as possible - but not simpler \
> -- <cite>Harry S. Plinkett</cite>

I've also heard that `tmux` could be used to automate this. But let's face it - if I'm too stupid to learn how to use `vim` then I assure you I'm way too stupid to learn `tmux`.

Last but not least - we're programmers for crying out loud. We get the unique luxury of creating the tools if they're not available. Because software development should be all about doing fun things, even if they turn out to be crap that nobody uses. Because it keeps our brain active, involved in a pure, creative process.

As I've said, this all started, practically, from a dare on a chat. I knew that docker-compose defined the syntax for what I wanted to achieve, it just didn't use the tooling that I liked. I thought - how hard would it be to create a tool that would fit my expectations?
