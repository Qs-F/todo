# package todo(second-version)

this package provides very simple todo managing.

## design

todo's thought is based git repository thought.  
when you add new todo, todo command searchs todo file.  
if there is no todo file on the same directory, then try to search on the upper directory.  
if there is no todo file on upper directory, create todo file on current directory. 

we call this "project". please remember it.

you can't register multiple lines todo. please write in one line form.  
(why?)  
(I want to create simple todo manager. if I do so, this will be not along with my original thought.)

## usage

### booked words

there are some disusable words for message.  
if you would like to add such words, you have to use -m option.

- init
- all
- archive

### details

- todo

this shows active todo order by registered date.

- todo all

this shows all todo(active or non-active)-

- todo MESSAGE

you can register new todo.-

- todo archive

this shows non-active(archived) todo.

- todo archive add

open choosing console and archive it.

- todo init

this creates new todo file on the same directory.-
