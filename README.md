# CopyThat
Go program that copies barcodes. 
Contains functions to use both Brady and Zebra printers. 
The printing is handled via command line and the PrintFile program created by Peter Lerup (available at https://www.lerup.com/mysoft.phtml).
For most installations the folder containing prfile32.exe must be added to the environments PATH. 
Also, a setting must be created in PrintFile called "Zebra" or "Brady" so PrintFile knows what printer to use.
The GUI is very simple and uses the WALK windows wrapper (available at https://github.com/lxn/walk).
