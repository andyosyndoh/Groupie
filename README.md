## Groupie-tracker-search-bar

### Introduction

* Groupie tracker search bar consists of creating a functional program that searches, inside groupietracker website, for a specific text input.
### Description

* Groupie-tracker is a go project which consists on receiving a given [API](https://groupietrackers.herokuapp.com/api) and manipulate the data contained in it, in order to create a site, displaying the information.


* The focus search bar project is to create a way for the client to search a member or artist or any other attribute in the data system made.

### Features


* The program handles at least these search cases :
    * artist/band name
    * members
    * locations
    * first album date
    * creation date
* The program also handle search input as case-insensitive.
* The search bar also have typing suggestions as you write.
  * The search bar identifies and displays in each suggestion an individual type of the search cases. (ex: Freddie Mercury -> member)
    * For example if you start writing "phil" it should appear as suggestions Phil Collins - member and Phil Collins - artist/band. This is just an example of a display.

### usage

  * To use the program, one needs to follow the steps below:

  1. Clone into the repository:
  ```bash
  git clone https://learn.zone01kisumu.ke/git/aosindo/groupie-tracker.git
  ```

  2. change directory to `groupie-tracker`

  ```bash
  cd groupie-tracker
  ```
  3. change directory to `cmd`

   ```bash
  cd cmd
  ```

4. run the command:
```bash
go run .
```
5. go to the web port given, @http://localhost:8080

6. Search from the search button.

## Authors

[aosindo](https://github.com/andyosyndoh)

[vandisi](https://github.com/Vinolia-E)