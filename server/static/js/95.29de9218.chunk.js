(this["webpackJsonppigsty-config"]=this["webpackJsonppigsty-config"]||[]).push([[95],{333:function(n,e,t){!function(n){"use strict";n.defineMode("sieve",(function(n){function e(n){for(var e={},t=n.split(" "),i=0;i<t.length;++i)e[t[i]]=!0;return e}var t=e("if elsif else stop require"),i=e("true false not"),r=n.indentUnit;function u(n,e){var r=n.next();if("/"==r&&n.eat("*"))return e.tokenize=l,l(n,e);if("#"===r)return n.skipToEnd(),"comment";if('"'==r)return e.tokenize=s(r),e.tokenize(n,e);if("("==r)return e._indent.push("("),e._indent.push("{"),null;if("{"===r)return e._indent.push("{"),null;if(")"==r&&(e._indent.pop(),e._indent.pop()),"}"===r)return e._indent.pop(),null;if(","==r)return null;if(";"==r)return null;if(/[{}\(\),;]/.test(r))return null;if(/\d/.test(r))return n.eatWhile(/[\d]/),n.eat(/[KkMmGg]/),"number";if(":"==r)return n.eatWhile(/[a-zA-Z_]/),n.eatWhile(/[a-zA-Z0-9_]/),"operator";n.eatWhile(/\w/);var u=n.current();return"text"==u&&n.eat(":")?(e.tokenize=o,"string"):t.propertyIsEnumerable(u)?"keyword":i.propertyIsEnumerable(u)?"atom":null}function o(n,e){return e._multiLineString=!0,n.sol()?("."==n.next()&&n.eol()&&(e._multiLineString=!1,e.tokenize=u),"string"):(n.eatSpace(),"#"==n.peek()?(n.skipToEnd(),"comment"):(n.skipToEnd(),"string"))}function l(n,e){for(var t,i=!1;null!=(t=n.next());){if(i&&"/"==t){e.tokenize=u;break}i="*"==t}return"comment"}function s(n){return function(e,t){for(var i,r=!1;null!=(i=e.next())&&(i!=n||r);)r=!r&&"\\"==i;return r||(t.tokenize=u),"string"}}return{startState:function(n){return{tokenize:u,baseIndent:n||0,_indent:[]}},token:function(n,e){return n.eatSpace()?null:(e.tokenize||u)(n,e)},indent:function(n,e){var t=n._indent.length;return e&&"}"==e[0]&&t--,t<0&&(t=0),t*r},electricChars:"}"}})),n.defineMIME("application/sieve","sieve")}(t(51))}}]);
//# sourceMappingURL=95.29de9218.chunk.js.map