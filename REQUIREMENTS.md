A web application for simple journaling. Single user, no authentication is required.

The tech stack should be:

- sqlite for the database
- go(lang) for the backend
- svelte5 and typescript for the frontend, communicates with the backend via REST
- the frontend is embedded in the backend binary, for portability.
- at least binaries for amd64
- Dockerfile for building a docker image (two stages)
- Semantic versioning
- CSS Framework: use Bootstrap
- Tooling: Makefile
- The application should respond on port 8080, no HTTPS needed
- Visual stile should be the standard Bootstrap one

# General

The main entity is an **entry**, which is a journal "day" record. Each entry, has a number of fields: for example, a rating of the day, what I ate at lunch, a general overview, where I was... No field is mandatory. The fields are:

- "POSITION_LON": longitude, where was I
- "POSITION_LAT": latitude, where was I
- "POSITION_NAME": text, location name resolved from coordinates
- "RATING": 1 to 10, a rating of the day
- "GENERAL": text, a general overview of the day
- "WORKING": state, whether I worked ("yes", "partial", "no", or empty)
- "MOOD": 1 to 10, a rating of my mood
- "MOOD_TXT": one-liner text, a short description of the rating above
- "LUNCH": text, what I ate for lunch
- "DINNER": text, what I ate for dinner
- "TV": text, what we watched on TV
- "SLEEP": 1 to 10, a rating of my my sleep the night before
- "SLEEP_TXT": one-liner text, a short description of the rating above

# Database

A file `scripts/schema.sql` must be generated with the SQL schema (`CREATE TABLE` and so on).

The tables are mainly two: an header with info on each journal entry, and a detail with the fields.

- [TABLE] `entries` -- header for an entry
  - [COLUMN] `date` TEXT -- date for the journal entry, choose a format that is optimal (text, datetime...)
  - [COLUMN] `timestamp` TEXT -- date and time, when the entry was last updated, choose a format that is optimal (text, datetime...)
  - [PRIMARY KEY] `date`

- [TABLE] `details` -- details for an entry
  - [COLUMN] `date` TEXT -- reference to `entries`.`date`
  - [COLUMN] `field_id` TEXT -- the ID of the field listed above ("GENERAL" ecc.)
  - [COLUMN] `value` TEXT -- the value of that ID, or NULL if it wasn't filled in
  - [PRIMARY KEY] `date` + `field_id`
  - [FOREIGN KEY] `date` -> `entries`.`date`

At least these tables must be created, but anything that is also necessary can be added. This way, in the future it will be easy to insert a new field.

Also, this structure will map nicely to JSON.

# UI

The UI should be optimized for desktop and for mobile (responsive).

The main UI view is a calendar view. It should mark the days that have an entry, and today. The current month should be shown by default, with arrows to go to the previous months.

Three buttons should be prominent, centered above the calendar:

- `View` to view the entry for the selected day
- `Edit` to create/edit an entry
- `Delete` will delete the selected entry, after a confirmation request

When the UI is loaded, the current day should be selected. It must be possible to select another day. The `Edit` button must be always active, whereas the `View` and `Delete` buttons must be available only if the selected day has an entry.

When pressing `View` or `Edit`, the same details form must be shown. The only difference is that the controls will be inactive if the `View` button was used.

This details form will be a form with the entry details, as described above in the General section, in that order. Some notes:

- `text` indicates a text area, `one-liner text` a text field, `1 to 10` a stars control to click, `state` a button group with options
- For the ratings (RATING, MOOD, SLEEP) the votes are 1-10, so 0 stars indicates "no rating". If "no rating" is given, the corresponding _TXT field (where present) should not be active. Also, the number of selected stars should be indicated beside the control, 'N.a.' if 0.
- Longitude and latitude must be implemented with a button "Current position" that hooks to the location system of the phone
- A button "Get location name" should reverse-geocode the coordinates to a location name using OpenStreetMap Nominatim API (no API key required)
- When a position (long/lat) is present, a link to Google Maps (external) should be given.
- WORKING field should be a button group with three options: "Worked" (green), "Partial" (yellow), "Not worked" (red). Only one can be selected at a time, or none.
- Calendar days should have different background colors based on WORKING status: green for worked, yellow for partial, red for not worked, blue for entries without working status.
- If a text area/field is left blank, or no rating, or no position, the corresponding detail should be NULL but the detail line must be present

2 buttons should be at the top of the form, centered:

- `Save` will save the form, and should be inactive if the user got here using `View`; 
- `Cancel` should go back to the main view without saving.

Make it beautiful, lean and fast. Code should be terse and well-commented.
