package apos

const(
	//font style
	COLOR_STYLE_RESET           = ";0"
	COLOR_STYLE_HIGHLIGHT       = ";1"
	COLOR_STYLE_DIM             =";2"
	COLOR_STYLE_UNDERLINE       =";4"
	COLOR_STYLE_REV             =";7"

	//front color
	COLOR_FRONT_BLACK           =";30"
	COLOR_FRONT_RED             =";31"
	COLOR_FRONT_GREEN           =";32"
	COLOR_FRONT_YELLOW          =";33"
	COLOR_FRONT_BLUE            =";34"
	COLOR_FRONT_PINK            =";35"
	COLOR_FRONT_CYAN            =";36"
	COLOR_FRONT_WHITE           =";37"

	//background color
	COLOR_BACK_BLACK           =";40"
	COLOR_BACK_RED             =";41"
	COLOR_BACK_GREEN           =";42"
	COLOR_BACK_YELLOW          =";43"
	COLOR_BACK_BLUE            =";44"
	COLOR_BACK_PINK            =";45"
	COLOR_BACK_CYAN            =";46"
	COLOR_BACK_WHITE           =";47"

	//prefix
	COLOR_PREFIX                    ="\033["
	COLOR_SUFFIX                    ="m"

	//short struct
	COLOR_SHORT_RESET         =COLOR_PREFIX+COLOR_STYLE_RESET+COLOR_SUFFIX
)

/*
HOW TO USE THE COLOR CONTROL:
LIKE THIS:
PREFIX + CMD + SUFFIX

EX1:RED
COLOR_PREFIX + COLOR_FRONT_RED + COLOR_SUFFIX

EX2:FRONT RED AND BACKGROUND GREEN
COLOR_PREFIX + COLOR_FRONT_RED + COLOR_BACK_GREEN + COLOR_SUFFIX

EX3:RESET TO THE DEFAULT
COLOR_SHORT_REST

NOTICE:
if you use fmt.println(PREFIX + CMD + SUFFIX , "This is a color show" , COLOR_SHORT_RESET)
The last show:
 This is a color show,(focus:before content "This is a color show" , a space here)

here you will see a space before Thisxxxxxx ,why? because the fmt.println(a,b,c),the println will add a space
in a and b ,so you will see a space before your content.

*/

