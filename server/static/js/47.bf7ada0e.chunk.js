(this["webpackJsonppigsty-config"]=this["webpackJsonppigsty-config"]||[]).push([[47],{279:function(t,e,a){!function(t){"use strict";t.defineMode("ebnf",(function(e){var a={slash:0,parenthesis:1},c={comment:0,_string:1,characterClass:2},r=null;return e.bracesMode&&(r=t.getMode(e,e.bracesMode)),{startState:function(){return{stringType:null,commentType:null,braced:0,lhs:!0,localState:null,stack:[],inDefinition:!1}},token:function(e,n){if(e){switch(0===n.stack.length&&('"'==e.peek()||"'"==e.peek()?(n.stringType=e.peek(),e.next(),n.stack.unshift(c._string)):e.match("/*")?(n.stack.unshift(c.comment),n.commentType=a.slash):e.match("(*")&&(n.stack.unshift(c.comment),n.commentType=a.parenthesis)),n.stack[0]){case c._string:for(;n.stack[0]===c._string&&!e.eol();)e.peek()===n.stringType?(e.next(),n.stack.shift()):"\\"===e.peek()?(e.next(),e.next()):e.match(/^.[^\\\"\']*/);return n.lhs?"property string":"string";case c.comment:for(;n.stack[0]===c.comment&&!e.eol();)n.commentType===a.slash&&e.match("*/")||n.commentType===a.parenthesis&&e.match("*)")?(n.stack.shift(),n.commentType=null):e.match(/^.[^\*]*/);return"comment";case c.characterClass:for(;n.stack[0]===c.characterClass&&!e.eol();)e.match(/^[^\]\\]+/)||e.match(".")||n.stack.shift();return"operator"}var s=e.peek();if(null!==r&&(n.braced||"{"===s)){null===n.localState&&(n.localState=t.startState(r));var i=r.token(e,n.localState),m=e.current();if(!i)for(var h=0;h<m.length;h++)"{"===m[h]?(0===n.braced&&(i="matchingbracket"),n.braced++):"}"===m[h]&&(n.braced--,0===n.braced&&(i="matchingbracket"));return i}switch(s){case"[":return e.next(),n.stack.unshift(c.characterClass),"bracket";case":":case"|":case";":return e.next(),"operator";case"%":if(e.match("%%"))return"header";if(e.match(/[%][A-Za-z]+/))return"keyword";if(e.match(/[%][}]/))return"matchingbracket";break;case"/":if(e.match(/[\/][A-Za-z]+/))return"keyword";case"\\":if(e.match(/[\][a-z]+/))return"string-2";case".":if(e.match("."))return"atom";case"*":case"-":case"+":case"^":if(e.match(s))return"atom";case"$":if(e.match("$$"))return"builtin";if(e.match(/[$][0-9]+/))return"variable-3";case"<":if(e.match(/<<[a-zA-Z_]+>>/))return"builtin"}return e.match("//")?(e.skipToEnd(),"comment"):e.match("return")?"operator":e.match(/^[a-zA-Z_][a-zA-Z0-9_]*/)?e.match(/(?=[\(.])/)?"variable":e.match(/(?=[\s\n]*[:=])/)?"def":"variable-2":-1!=["[","]","(",")"].indexOf(e.peek())?(e.next(),"bracket"):(e.eatSpace()||e.next(),null)}}}})),t.defineMIME("text/x-ebnf","ebnf")}(a(51))}}]);
//# sourceMappingURL=47.bf7ada0e.chunk.js.map