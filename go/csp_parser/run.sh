# install antlr
export CLASSPATH=".:/usr/local/lib/antlr-4.9-complete.jar:$CLASSPATH"
alias antlr4='java -jar /usr/local/lib/antlr-4.9-complete.jar'
alias grun='java org.antlr.v4.gui.TestRig'

# compile to runable
antlr4 Csp.g4 -visitor
javac *.java
grun Csp spec -gui <filename>.csp