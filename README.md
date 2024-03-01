# habittui
This is my terminal application for managing habits and maybe other recurring tasks.

## Initial project:
The terminal will be split in at least 3 separate windows. Name window is used interchangeably as section.
```
----------------------------------------------------------|
| Tasks (today):             | Description:               |
|  (short name)              |                            |
|  Work on habitui   [x]     |  Extended description of   |
|  Go for a walk     [x]     |  task that can be edited   |
|  Do english lesson [ ]     |                            |
|                            |                            |
| Shows tasks for today.     |                            |
|----------------------------|-----------------------------
| Statistics: (of selected   | To be decided what         |
|              task)         |   goes in this section.    |
|  Strike: 2 days in a row.  |                            |
|  This week: 3 times        |                            |
|  This month: 10 times      |                            |
|  This year: 30 times       |                            |
|----------------------------|----------------------------|
```
## Navigation
j - Navigate up in current window, in Tasks windows this will highlight (or change color of) currently selected task. <br>
k - Navigate down in current windows, in Tasks windows. <br>
l - Navigate to window right, for example from tasks to description window. This will somehow highlight (or color change) other window.  <br>
h - Navigate to window left. <br>
e - Edit data in current window. Allows to edit task short name or description.
        After edit, user press enter and a pop up with confirm changes 'y/n' will come. <br>
a/c - To be decided if both (add/create) - creates new habit task. Will be based on popup/floating windows that ask you questions
      like: enter (short name), enter long description. <br>
      Alternatively this will make clear terminal window and open a new one with empty task short name and description.
      User will move over short name and description using hl keys press 'i' (insert) to add data.  <br>
d - Deletes currently selected task, a confirmation window will pop up. <br>
s - Saves changes that were made during current session (to be decided if this will be actually handled or saves will be done
automatically).
r - Restart session (revert), reverts all (or last) changes that were done during 
p - period (to be decided), changes period that we want to manage (by default daily), but some tasks could be defined to be 
monthly, yearly etc. for example monthly task pay this month bills.
## Consistent state 
The app should keep changes that were save in the initial (stage), possibly each change should be kept even if app is restarted
unless declined. Initial implementation will possibly keep changes in local text file like json or csv. Later it might be moved
to some database.
## Configurability
To be decided what mechanism should be changed and possibly modifiable using command line args or config file.

