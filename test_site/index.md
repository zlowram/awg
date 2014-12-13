Welcome to AWG, dear user!

If you are reading this it means you already know how to use the CLI program in
order to generate a site. Please, read further in order to know how to create a
site and the basic syntax.

The config file
---------------

The configuration for each site must be defined in the awg.conf file. The
config file for this site is the following:

--- awg.conf -----------------------

{
    "Title":    "awg docs",
    "Logo":     "logo.txt",
    "Style":    "themes/dark.css"
}

-------------------------------------

Title, as the name itself says, defines the title of the site to be generated.
This will appear between the "title" tags in the header of the resulting HTML
file.

Logo, defines the file where the ASCII logo of the site is stored.

Style, defines the stylesheet to be used on the generated site. You can find 3
basic styles in the theme folder. 

The menu
--------

The menu you see below the site logo is automatically generated following the
directory structure of your site. For example, this test site is organized as
follows:

  + index.md
  | 
  | 
  + syntax/
  |   | 
  |   +-- index.md
  |   | 
  |   +-- images.md
  |   | 
  |   +-- links.md
  | 
  + about.md
  
Now that you know the directory structure, go and check the menu!
