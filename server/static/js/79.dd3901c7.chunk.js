(this["webpackJsonppigsty-config"]=this["webpackJsonppigsty-config"]||[]).push([[79],{318:function(t,e,n){!function(t){"use strict";t.defineMode("pegjs",(function(e){var n=t.getMode(e,"javascript");function r(t){return t.match(/^[a-zA-Z_][a-zA-Z0-9_]*/)}return{startState:function(){return{inString:!1,stringType:null,inComment:!1,inCharacterClass:!1,braced:0,lhs:!0,localState:null}},token:function(e,i){if(e&&(i.inString||i.inComment||'"'!=e.peek()&&"'"!=e.peek()||(i.stringType=e.peek(),e.next(),i.inString=!0)),i.inString||i.inComment||!e.match("/*")||(i.inComment=!0),i.inString){for(;i.inString&&!e.eol();)e.peek()===i.stringType?(e.next(),i.inString=!1):"\\"===e.peek()?(e.next(),e.next()):e.match(/^.[^\\\"\']*/);return i.lhs?"property string":"string"}if(i.inComment){for(;i.inComment&&!e.eol();)e.match("*/")?i.inComment=!1:e.match(/^.[^\*]*/);return"comment"}if(i.inCharacterClass)for(;i.inCharacterClass&&!e.eol();)e.match(/^[^\]\\]+/)||e.match(/^\\./)||(i.inCharacterClass=!1);else{if("["===e.peek())return e.next(),i.inCharacterClass=!0,"bracket";if(e.match("//"))return e.skipToEnd(),"comment";if(i.braced||"{"===e.peek()){null===i.localState&&(i.localState=t.startState(n));var a=n.token(e,i.localState),c=e.current();if(!a)for(var o=0;o<c.length;o++)"{"===c[o]?i.braced++:"}"===c[o]&&i.braced--;return a}if(r(e))return":"===e.peek()?"variable":"variable-2";if(-1!=["[","]","(",")"].indexOf(e.peek()))return e.next(),"bracket";e.eatSpace()||e.next()}return null}}}),"javascript")}(n(50),n(249))}}]);
//# sourceMappingURL=79.dd3901c7.chunk.js.map