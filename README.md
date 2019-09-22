# Production Grid Core

This project provides the core foundational elements for rapid business application
development in Go.  Production Grid provides a base of elements almost every web
application needs like security and user management so developers don't have to
waste time building things they've already built a million times before on other
projects.

Developers extend Production Grid by adding modules.  Modules can add permissions
to the security infrastructure, events, API's, server rendered pages, database tables,
etc.

The core library provides the following foundational elements:

* A Common Sense DAO Abstraction
* Database Schema Migration Tools
* Reporting and PDF Rendering Tools
* User Management with Password Resets and Security Groups
* Session Management
* A notification system with support for SMS and EMail transport types.
* Ayncronous Processing Support (via Channels or message queues)
* Web Template Management
* Event Processing and Alerting
* Symmetric Encryption
* Scheduled Event and Jobs
* Bundled resources, API, and server rendered content pages.
