(this["webpackJsonppigsty-config"]=this["webpackJsonppigsty-config"]||[]).push([[44],{274:function(i,n,e){!function(i){"use strict";i.defineMode("diff",(function(){var i={"+":"positive","-":"negative","@":"meta"};return{token:function(n){var e=n.string.search(/[\t ]+?$/);if(!n.sol()||0===e)return n.skipToEnd(),("error "+(i[n.string.charAt(0)]||"")).replace(/ $/,"");var t=i[n.peek()]||n.skipToEnd();return-1===e?n.skipToEnd():n.pos=e,t}}})),i.defineMIME("text/x-diff","diff")}(e(51))}}]);
//# sourceMappingURL=44.62a7f353.chunk.js.map