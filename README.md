# habittui
This is my terminal application for managing habits and maybe other recurring tasks.
Rationale for this app is that I spend most of the day in terminal as my setup is neo-vim + tmux. 
Having habit tracking app inside of terminal is more convenient for me than mobile or web app. 
## Initial project:
The terminal will be split in at least 3 separate windows. Name window is used interchangeably as section.
```
----------------------------------------------------------|
| Tasks (today):             | Description:               |
|  (short name)              |                            |
| [x] Work on habitui        |  Extended description of   |
| [x] Go for a walk          |  task that can be edited   |
| [ ] Do english lesson      |                            |
|                            |                            |
| Shows tasks for today.     |                            |
|----------------------------|-----------------------------
| Strike statistics:         | Completion statistics:     |
|                            |                            |
|   Current: 2 days          | This week: 2 times         |
|   Best monthly: 5 days     | This month: 10 times       |
|   Longest: 10 days         | This year: 50 times        |
|                            |                            |
|                            |                            |
|----------------------------|----------------------------|
```
## Navigation
j - Navigate up in current window, in Tasks windows this will highlight (or change color of) currently selected task. <br>
k - Navigate down in current windows, in Tasks windows. <br>
l - Navigate to window right, for example from tasks to description window. This will somehow highlight (or color change) other window.  <br>
h - Navigate to window left. <br>
e - Edit data in current window. Allows to edit task short name or description.
        After edit, user press enter and a pop up with confirm changes 'y/n' will come. <br>
a - add new task and move into task name/description edit mode. <br>
d - Deletes currently selected task, a confirmation window will pop up. <br>
p - period (to be decided), changes period that we want to manage (by default daily), but some tasks could be defined to be 
monthly, yearly etc. for example monthly task pay this month bills. <br>
## Consistent state 
The app should keep changes that were save in the initial (stage), possibly each change should be kept even if app is restarted
unless declined. Initial implementation will possibly keep changes in local text file like json or csv. Later it might be moved
to some database.
## Configurability
To be decided what mechanism should be changed and possibly modifiable using command line args or config file.

## Developement status
App dev in progress.
