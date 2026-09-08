## Purpose

Defines real-time type-to-filter behavior for the interactive EC2 instance picker, so users can narrow a long instance list by typing instead of scrolling, consistently across every command that presents the picker.

## ADDED Requirements

### Requirement: Interactive instance picker supports type-to-filter
The interactive instance selector SHALL allow the user to filter the displayed instances in real time by typing a query. The filter SHALL be a case-insensitive substring match against each instance's `InstanceID`, `Name`, `PrivateIP`, and `PublicIP` fields, consistent with the positional `[search]` fragment matching used by `ec2connect`. Only matching instances SHALL remain selectable while a filter is active. This behavior SHALL apply to every command that presents the instance picker (`ec2connect`, `ec2`, `ec2stopstart`, `ec2scheduler`).

#### Scenario: Typing narrows the list in real time
- **GIVEN** the interactive instance picker is showing multiple instances
- **WHEN** the user types a query fragment
- **THEN** the picker SHALL show only instances whose identifying fields contain the fragment (case-insensitive)
- **AND** SHALL update as each character is typed

#### Scenario: Selecting a filtered instance
- **GIVEN** a filter is active and one or more instances match
- **WHEN** the user moves the highlight to a matching instance and confirms the selection
- **THEN** the system SHALL return that instance and proceed as if it had been selected from the unfiltered list

#### Scenario: Clearing the filter restores the full list
- **GIVEN** a filter is active
- **WHEN** the user clears the filter query
- **THEN** the picker SHALL again show all instances

#### Scenario: No matches
- **GIVEN** the interactive instance picker is showing instances
- **WHEN** the user types a query that matches no instance
- **THEN** the picker SHALL show no selectable instances
- **AND** SHALL allow the user to edit or clear the query to continue

#### Scenario: Cancel semantics preserved
- **WHEN** the user cancels the picker (Esc or Ctrl+C), whether or not a filter is active
- **THEN** the system SHALL return no selection and SHALL NOT connect to or act on any instance

#### Scenario: Consistent across commands using the picker
- **WHEN** the instance picker is shown by `ec2connect`, `ec2`, `ec2stopstart`, or `ec2scheduler`
- **THEN** the same type-to-filter behavior SHALL be available in each
