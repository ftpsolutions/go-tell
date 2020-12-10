### go-tell ###

#### What is this? ####
Go-tell is a set of opionated interfaces and structs intended to be used as building blocks for a notification system.

The blocks are intended to be constructed as so:

    Source -> Store <--> Worker -> Sender

##### Source
 - A source knows how to receive a job and save to a store

##### Store
 A store is made up of two parts;
 - Store is the behavioural operations tailored specifically to managing jobs.
 - Storage is the 'CRUD' operations used for jobs.

##### Worker
 - A worker knows how to 'handle' a job from a store to a specific sender interface. 
 - A worker has a 'strategy' to manage jobs if they fail

##### Sender
 - A sender tries to send and stops on failure.
 - Inspired by the sender interfaces at 'github.com/appscode/go-notify', there's a wrapper for them.

 #### So what can I even with this? ####
 For now, there's a lot of plumbing required. This is only a library with a shared set of interfaces that -you- need to put together.

 Check the examples for simple implementations.

 #### TODO ####
  - Add sources
  - More default job handlers
  - Pass error handling behaviour through to a worker
  - Have a default set of error handling behaviours for workers
  - Option to have environment configurations as defaults for stores/senders/sources
  - A full system (cmd/main.go)
