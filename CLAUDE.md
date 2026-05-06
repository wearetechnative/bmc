- always remember openspec settings and documents
- never add claude or claude-code references in git commit messages
When I refer to issues like isotto-rn3b checkout the task
in @.beans/isotto-rn3b*.md
In this project we will use these tasks as epics for making openspec proposals.
WHEN you create a proposal at a link to this task in the proposal.md.
WHEN a bean is used to create an proposal change the status to "in-progress"
WHEN a proposal is archived add the link to the archived proposal in the frontmatter of this task like this:
openspec-link: openspec/changes/archive/....
You are allowed to update these statuses in the task frontmatter:
in-progress
todo
draft
completed
scrapped
When making changes you are allowed to update the date/time in updated_at in the task frontmatter
Besides updating status and openspec-link, you are NOT ALLOWED to modify the contents of the task file.
Always use opsx commands when creating openspec proposals or archive proposals
All openspec documents need to be in english, not matter the language being used in the users conversation.

