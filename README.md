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

* A Common Sense DAO Abstraction based on Postgres.
* Database Schema Migration Tools with pre and post migration phase for zero downtime deployments.
* Reporting and PDF Rendering Tools
* User Management with Password Resets and Security Groups
* Multi-Tenant Support with Application Definable Terminology
* Session Management
* A notification system with support for SMS and EMail transport types.
* Ayncronous Processing Support (via Channels or message queues)
* Web Template Management
* Event Processing and Alerting
* Media Archival, Scaling, and Transcoding
* Symmetric Encryption
* Scheduled Event and Jobs
* Bundled resources, API, and server rendered content pages.

## Roadmap Notes

This library will provide core content management, e-commerce, and application
infrastructure.  When complete, it will essentially be a full content management
system with e-commerce built into the base version.

The primary purpose behind this project is to support a set of closed source commercial
e-commerce and production management tools for Performing Arts Organizations while leaving
the parts with more broad utility in the public domain.

The commercial projects will extend the core library with the following features:

* Cast and Crew Scheduling (with Mobile Apps)
* Artistic and Work Notes Distribution/Archival
* Production Collaboration Tools (Document Sharing, Chat)
* Child Protection Features for Shows Involving Young Performers
* Integrated Background Checks for Adults Working With Young Performers
* Performance Reports
* Digital Sign In Sheet with PSM Notifiations
* Wiki Style Public Database of Equipment and Venue Specifications
* Equipment Inventory
* Lighting Data Management Tools (i.e. Cloud based Lightwright)
* Production Budgeting and Expense Tracking
* Online Ticketing and Box Office Point-of-Sale (with chip card and ApplePay support)
* Ingress Ticket Scanning
* BOCA Printer Integration
* Donor Tracking and Fundraising Event Tools (e.i. End of Night Settlement)
* One Time Donations and Recurring Donations
* Class Registration and Scheduling
* Tax Letter Generation
* Integrated Royalty and Performance Rights Management
