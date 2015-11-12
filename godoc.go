/*voice

A fast, distributed, scalable bulletin board (forum) written in Go using event sourcing.

Modules (packages)

The application is broken into what can be considered a "micro-services" architecture, which fits nicely with DDD's concept of Modules.  These modules are organized as golang packages in this application.  Modules implement behaviors, grouping related functionality in the same package and hiding dirty implementation details behind the nice contract. The contract is expressed as an interface:

- Modules can subscribe to events (user approved, photo liked, payment received etc)
- Modules can expose HTTP REST endpoints (approve photo, register, get account details etc)
- Modules can publish events

Modules can have distinct data storages. They interact only through the interface (and can be deployed on different machines).

We don't have any notion of aggregates, aggregate roots, repositories or process managers in the design, they weren't needed. We might have these concepts inside the module, but that would be an irrelevant implementation detail. Event sourcing is aligned with the definition of Martin Fowler (not the way it is usually done in .NET/Java communities).

Modules are relatively small, golang allows that. They merely group related behaviors together. In that sense their design process is similar to how one would group similar behavior and state in aggregates. A bounded context can be implemented as a dozen of small modules, if a language permits that easily. Golang does that.

Modules are constrained to be a separate unit of deployment and scaling (we can deploy them to different machines in multiple instances).

To be more precise, domain design, expressed in modules, their contracts and relations is important. Implementation details are less relevant, even the language in which they are written. If you have a decent domain design, throwing out the code and rewriting modules one by one is trivial.

Although, not much emphasis was placed on the implementation details, there were a few guidelines:

- We denormalize heavily in the modules. Module could subscribe to the needed events and maintain a local view model, used for decision making, enriching HTTP responses or published events. It works ok for systems with millions of events and above.
- Modules are idempotent. This is verified by running tests, derived from use-cases.
- Modules tend to be small - a few tables, a few queries, a few HTTP routes.
- Performance in this project is so important that it became a part of the domain (we have big shoes to fill). This added additional design constraints to the process, making it more technical.

Naming the module was probably the hardest thing (especially, since package names in golang tend to be short).

To refine the definition of modules even further:

- our application is composed from modules - tangible way to structure code and visual way to group design concepts into;
- we align modules at the design level and modules in the code - they share the same boundaries and names;
- at design level modules are boxes that have associated behavior, we need them to contain complexity and decompose our design into small concepts that are easy to reason and talk about;
- in the codebase our modules are represented by folders which also are treated as packages and namespaces in golang;
- we like to keep our modules small, focused and decoupled, this requires some discipline but speeds up development;
- each module has its own public contract by which it is known to the other modules; implementation details are private, they can't be coupled to and are treated as black-box;
- Public contract can include: published events (events are a part of domain language), public golang service interfaces and http endpoints; there also are behavioral contracts setting up expectations on how these work together;
- in the code each golang package is supposed to have an implementation of the following interface, that's how it is wired to the system; all module dependencies are passed into the constructor without any magic.

  type Module interface {
    Register(c Context)
  }

  type Context interface {
    AddAuthHttp(path string, handler web.Handler)
    AddHttpHandler(path string, handler http.Handler)
    RegisterEventHandler(h bus.NodeHandler)
    ResetData(reset func())
  }

Getting the notion of modules right is extremely important for us, since it is one of the principles behind our design process. We think, structure our work and plan in terms of modules.

The modules will be listening to the event streams.  Therefore, workflow through the system for testing is as simple as:

  func run_nancy_flirts_bob(x *context) (info *nancy_flirts_bob) {
    info = &nancy_flirts_bob{voice.NewId(), voice.NewId()}
    x.Dispatch(voice.NewRegistrationApproved(
      voice.NewId(),
      info.bobId,
      "bob",
      voice.Male,
      voice.NewBirthday(time.Now().AddDate(-23, 0, 0)),
      "email",
      voice.NoPortraitMale))

    x.Dispatch(voice.NewRegistrationApproved(
      voice.NewId(),
      info.nancyId,
      "nancy",
      voice.Female,
      voice.NewBirthday(time.Now().AddDate(-22, 0, 0)),
      "email",
      voice.NoPortraitFemale))

    x.Dispatch(&voice.FlirtSent{voice.NewId(), info.nancyId, info.bobId})
    return info
  }

Data Storage

The application defaults to an in-memory service bus for event sourcing, making it easier to run standalone during development.

The production data adapter utilized is the EventStore (http://geteventstore.com) due to its replication and distributed nature.  The EventStore requires decent disk IOPS for high performance.  If writes become a bottleneck, considering upgrading to SSD drives and/or setup a replication cluster.

*/
package main
