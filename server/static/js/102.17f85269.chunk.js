(this["webpackJsonppigsty-config"]=this["webpackJsonppigsty-config"]||[]).push([[102],{256:function(t,e,n){!function(t){"use strict";t.defineMode("stex",(function(t,e){function n(t,e){t.cmdState.push(e)}function r(t){return t.cmdState.length>0?t.cmdState[t.cmdState.length-1]:null}function a(t){var e=t.cmdState.pop();e&&e.closeBracket()}function i(t){for(var e=t.cmdState,n=e.length-1;n>=0;n--){var r=e[n];if("DEFAULT"!=r.name)return r}return{styleIdentifier:function(){return null}}}function c(t,e,n){return function(){this.name=t,this.bracketNo=0,this.style=e,this.styles=n,this.argument=null,this.styleIdentifier=function(){return this.styles[this.bracketNo-1]||null},this.openBracket=function(){return this.bracketNo++,"bracket"},this.closeBracket=function(){}}}var u={};function o(t,e){t.f=e}function f(t,e){var a;if(t.match(/^\\[a-zA-Z@]+/)){var c=t.current().slice(1);return n(e,a=new(a=u.hasOwnProperty(c)?u[c]:u.DEFAULT)),o(e,s),a.style}if(t.match(/^\\[$&%#{}_]/))return"tag";if(t.match(/^\\[,;!\/\\]/))return"tag";if(t.match("\\["))return o(e,(function(t,e){return m(t,e,"\\]")})),"keyword";if(t.match("\\("))return o(e,(function(t,e){return m(t,e,"\\)")})),"keyword";if(t.match("$$"))return o(e,(function(t,e){return m(t,e,"$$")})),"keyword";if(t.match("$"))return o(e,(function(t,e){return m(t,e,"$")})),"keyword";var f=t.next();return"%"==f?(t.skipToEnd(),"comment"):"}"==f||"]"==f?(a=r(e))?(a.closeBracket(f),o(e,s),"bracket"):"error":"{"==f||"["==f?(n(e,a=new(a=u.DEFAULT)),"bracket"):/\d/.test(f)?(t.eatWhile(/[\w.%]/),"atom"):(t.eatWhile(/[\w\-_]/),"begin"==(a=i(e)).name&&(a.argument=t.current()),a.styleIdentifier())}function m(t,e,n){if(t.eatSpace())return null;if(n&&t.match(n))return o(e,f),"keyword";if(t.match(/^\\[a-zA-Z@]+/))return"tag";if(t.match(/^[a-zA-Z]+/))return"variable-2";if(t.match(/^\\[$&%#{}_]/))return"tag";if(t.match(/^\\[,;!\/]/))return"tag";if(t.match(/^[\^_&]/))return"tag";if(t.match(/^[+\-<>|=,\/@!*:;'"`~#?]/))return null;if(t.match(/^(\d+\.\d*|\d*\.\d+|\d+)/))return"number";var r=t.next();return"{"==r||"}"==r||"["==r||"]"==r||"("==r||")"==r?"bracket":"%"==r?(t.skipToEnd(),"comment"):"error"}function s(t,e){var n=t.peek();return"{"==n||"["==n?(r(e).openBracket(n),t.eat(n),o(e,f),"bracket"):/[ \t\r]/.test(n)?(t.eat(n),null):(o(e,f),a(e),f(t,e))}return u.importmodule=c("importmodule","tag",["string","builtin"]),u.documentclass=c("documentclass","tag",["","atom"]),u.usepackage=c("usepackage","tag",["atom"]),u.begin=c("begin","tag",["atom"]),u.end=c("end","tag",["atom"]),u.label=c("label","tag",["atom"]),u.ref=c("ref","tag",["atom"]),u.eqref=c("eqref","tag",["atom"]),u.cite=c("cite","tag",["atom"]),u.bibitem=c("bibitem","tag",["atom"]),u.Bibitem=c("Bibitem","tag",["atom"]),u.RBibitem=c("RBibitem","tag",["atom"]),u.DEFAULT=function(){this.name="DEFAULT",this.style="tag",this.styleIdentifier=this.openBracket=this.closeBracket=function(){}},{startState:function(){return{cmdState:[],f:e.inMathMode?function(t,e){return m(t,e)}:f}},copyState:function(t){return{cmdState:t.cmdState.slice(),f:t.f}},token:function(t,e){return e.f(t,e)},blankLine:function(t){t.f=f,t.cmdState.length=0},lineComment:"%"}})),t.defineMIME("text/x-stex","stex"),t.defineMIME("text/x-latex","stex")}(n(50))}}]);
//# sourceMappingURL=102.17f85269.chunk.js.map