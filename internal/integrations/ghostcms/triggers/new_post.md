🧭 Overview

This trigger queries the Ghost Admin API for posts, optionally filtered by status and with related authors/tags included.
It uses the workflow runtime’s lastRun timestamp to fetch only posts created after the previous execution. Results are returned oldest → newest so downstream steps process them in order.

⚙️ Properties (Inputs)
Field	Type	Required	Default	Description
status	string	No	all	Filter posts by status: all, published, draft, scheduled.
include	string	No	authors,tags	Comma-separated relations to include, e.g. authors,tags.