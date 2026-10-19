# JobTrack — Product Requirements Document (PRD)

**Version:** 1.0  
**Project duration:** 7 days  
**Status:** Proposed for team approval  
**Product type:** Web application  
**Primary users:** Job seekers

## 1. Product Overview

JobTrack is a web application that helps job seekers organize job applications, monitor application deadlines, track hiring progress, and compare their CV text against job requirements using keyword matching.

The application will provide a central dashboard where users can view their applications, identify approaching deadlines, and understand which keywords from a job description appear in their CV.

The first release will prioritize reliability, usability, and core functionality over advanced automation.

## 2. Problem Statement

Job seekers often manage multiple applications across different companies and job boards. Keeping track of application dates, deadlines, interview progress, and job requirements can become difficult when information is spread across spreadsheets, emails, and notes.

JobTrack aims to simplify this process by bringing essential job-search management features into one application.

## 3. Product Objectives

The MVP must:

1. Allow users to create, view, update, and delete job applications.
2. Track application status and relevant dates.
3. Identify upcoming and overdue application deadlines.
4. Compare CV text against job requirements using keyword matching.
5. Display useful application statistics on a dashboard.
6. Persist application data in a database.
7. Provide a responsive, understandable user interface.
8. Be easy for another developer to install, run, and test.

## 4. Target User

The primary user is a job seeker applying for internships, entry-level positions, graduate roles, or experienced professional positions.

The application will initially support a single local user environment without account registration or authentication. Multi-user accounts are outside the MVP scope.

## 5. Technology Stack

- **Frontend:** React, TypeScript, Vite
- **Backend:** Go and the standard `net/http` package
- **Database:** SQLite
- **API style:** REST
- **Data format:** JSON
- **Version control:** Git and GitHub
- **Testing:** Go unit and HTTP handler tests, plus frontend tests and/or documented manual test cases

The team must use the repository's agreed directory structure and document all required setup commands.

## 6. Functional Requirements

### FR-01: Create a Job Application

The user must be able to add an application with:

- Company name
- Job title
- Job URL (optional)
- Date applied (optional)
- Application deadline (optional)
- Application status
- Notes (optional)

**Acceptance criteria:**

- Company and job title are required.
- Valid submissions are saved to SQLite.
- The new application appears in the application list.
- Invalid submissions display a clear error.
- The record remains available after refreshing the page or restarting the backend.

### FR-02: View Job Applications

The user must be able to view all saved applications.

**Acceptance criteria:**

- Each record displays its company, job title, status, and relevant dates.
- An empty state is displayed when there are no applications.
- The user can open or inspect an individual application.
- Data displayed by the frontend reflects the API response.

### FR-03: Update an Application

The user must be able to update application details and status.

**Acceptance criteria:**

- Changes are saved to the database.
- The updated record appears in the interface.
- Invalid updates are rejected with a meaningful error.
- Updating a nonexistent record returns a not-found response.

### FR-04: Delete an Application

The user must be able to delete an application.

**Acceptance criteria:**

- The interface requests confirmation before deletion.
- A successful deletion removes the record from the list.
- Deleting a nonexistent record returns a not-found response.

### FR-05: Track Application Status

The application must support these statuses:

- `Saved`
- `Applied`
- `Interview`
- `Offer`
- `Rejected`

**Acceptance criteria:**

- Every application has exactly one valid status.
- New applications default to `Saved` when no status is supplied.
- The user can change status.
- Invalid status values are rejected.

### FR-06: Track Deadlines

The application must identify upcoming and overdue deadlines.

For the MVP:

- An application is **overdue** when its deadline is before the current calendar date and its status is not `Offer` or `Rejected`.
- An application is **upcoming** when its deadline falls between today and the next seven calendar days, inclusive, and its status is not `Offer` or `Rejected`.
- Applications without a deadline must not be marked overdue.
- Deadline comparisons must use calendar dates in `YYYY-MM-DD` format.

The backend is responsible for calculating deadline classifications returned by the API. The frontend may use those classifications to display labels and styling.

### FR-07: Keyword Matching

The user must be able to provide CV text and a job description for comparison.

The result must display:

- Match percentage
- Matched keywords
- Missing keywords
- Number of required keywords considered

The MVP will use deterministic keyword matching, not an AI model.

The matching implementation must document how keywords are identified. A predefined technology and skill vocabulary is acceptable for the first release. The user may also supply an explicit list of required keywords.

Matching must be case-insensitive and should avoid matching a keyword inside an unrelated longer word. Duplicate required keywords must not inflate the score.

If no keywords can be identified, the API must return a clear validation error rather than an arbitrary score.

### FR-08: Dashboard Statistics

The dashboard must show:

- Total applications
- Applications by status
- Upcoming deadlines
- Overdue applications
- Total interviews
- Total offers

Statistics must be calculated from persisted application data.

The definition of an interview is an application whose current status is `Interview`. The total offers count is the number of applications currently marked `Offer`.

### FR-09: Error Handling

The application must handle:

- Invalid request bodies
- Missing required fields
- Invalid dates
- Invalid status values
- Nonexistent application IDs
- Database errors
- Invalid keyword-matching requests
- Unexpected server errors

The API must return structured JSON errors and appropriate HTTP status codes. Internal database details and stack traces must not be exposed to the client.

## 7. Non-Functional Requirements

### Usability

- The interface must have clear navigation and readable labels.
- Forms must provide validation feedback.
- Empty, loading, success, and error states must be handled.
- The layout must remain usable on desktop and mobile screens.

### Reliability

- Application records must persist in SQLite.
- API handlers must handle invalid requests without crashing the server.
- The application must start using documented commands.
- All core workflows must be tested before submission.

### Security

- Validate user-supplied input on the backend.
- Use parameterized SQL queries.
- Do not expose database errors, secrets, or internal implementation details.
- Do not commit credentials or environment secrets.
- Treat CV text as potentially sensitive; do not log full CV contents.
- Do not render user-supplied text as executable HTML.

### Maintainability

- Keep API handlers, database access, and business logic separated where practical.
- Use clear names and consistent error handling.
- Document API routes and data structures.
- Keep functions small enough for beginner developers to understand and test.

## 8. User Stories

1. As a job seeker, I want to save a job so that I can keep its details in one place.
2. As a job seeker, I want to update an application status so that I can track my progress.
3. As a job seeker, I want to see approaching deadlines so that I can prioritize applications.
4. As a job seeker, I want to compare my CV with a job description so that I can identify missing keywords.
5. As a job seeker, I want to see application statistics so that I can understand my job-search progress.
6. As a job seeker, I want to edit or delete records so that my tracker remains accurate.


## 9. Definition of Done

A feature is considered complete when:

1. It satisfies its acceptance criteria.
2. It follows the agreed API contract.
3. Its relevant tests pass.
4. Invalid inputs and expected error conditions are handled.
5. It has been reviewed through a GitHub pull request.
6. It works with the integrated application.
7. Its setup or usage instructions are documented when necessary.


## 10. Success Criteria

The MVP is successful when:

- All eight functional requirements are implemented to the agreed scope.
- Application data persists across refreshes and backend restarts.
- Deadline classifications are correct for tested date boundaries.
- Keyword matching returns a consistent score and matched/missing keyword lists.
- Dashboard statistics agree with the saved records.
- The main user journey completes without critical defects.
- A team member who did not build the application can follow the README to run it.
