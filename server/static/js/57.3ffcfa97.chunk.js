(this["webpackJsonppigsty-config"]=this["webpackJsonppigsty-config"]||[]).push([[57],{291:function(t,e,n){!function(t){"use strict";t.defineMode("go",(function(e){var n,r=e.indentUnit,i={break:!0,case:!0,chan:!0,const:!0,continue:!0,default:!0,defer:!0,else:!0,fallthrough:!0,for:!0,func:!0,go:!0,goto:!0,if:!0,import:!0,interface:!0,map:!0,package:!0,range:!0,return:!0,select:!0,struct:!0,switch:!0,type:!0,var:!0,bool:!0,byte:!0,complex64:!0,complex128:!0,float32:!0,float64:!0,int8:!0,int16:!0,int32:!0,int64:!0,string:!0,uint8:!0,uint16:!0,uint32:!0,uint64:!0,int:!0,uint:!0,uintptr:!0,error:!0,rune:!0},o={true:!0,false:!0,iota:!0,nil:!0,append:!0,cap:!0,close:!0,complex:!0,copy:!0,delete:!0,imag:!0,len:!0,make:!0,new:!0,panic:!0,print:!0,println:!0,real:!0,recover:!0},a=/[+\-*&^%:=<>!|\/]/;function c(t,e){var r=t.next();if('"'==r||"'"==r||"`"==r)return e.tokenize=u(r),e.tokenize(t,e);if(/[\d\.]/.test(r))return"."==r?t.match(/^[0-9]+([eE][\-+]?[0-9]+)?/):"0"==r?t.match(/^[xX][0-9a-fA-F]+/)||t.match(/^0[0-7]+/):t.match(/^[0-9]*\.?[0-9]*([eE][\-+]?[0-9]+)?/),"number";if(/[\[\]{}\(\),;\:\.]/.test(r))return n=r,null;if("/"==r){if(t.eat("*"))return e.tokenize=l,l(t,e);if(t.eat("/"))return t.skipToEnd(),"comment"}if(a.test(r))return t.eatWhile(a),"operator";t.eatWhile(/[\w\$_\xa1-\uffff]/);var c=t.current();return i.propertyIsEnumerable(c)?("case"!=c&&"default"!=c||(n="case"),"keyword"):o.propertyIsEnumerable(c)?"atom":"variable"}function u(t){return function(e,n){for(var r,i=!1,o=!1;null!=(r=e.next());){if(r==t&&!i){o=!0;break}i=!i&&"`"!=t&&"\\"==r}return(o||!i&&"`"!=t)&&(n.tokenize=c),"string"}}function l(t,e){for(var n,r=!1;n=t.next();){if("/"==n&&r){e.tokenize=c;break}r="*"==n}return"comment"}function f(t,e,n,r,i){this.indented=t,this.column=e,this.type=n,this.align=r,this.prev=i}function s(t,e,n){return t.context=new f(t.indented,e,n,null,t.context)}function p(t){if(t.context.prev){var e=t.context.type;return")"!=e&&"]"!=e&&"}"!=e||(t.indented=t.context.indented),t.context=t.context.prev}}return{startState:function(t){return{tokenize:null,context:new f((t||0)-r,0,"top",!1),indented:0,startOfLine:!0}},token:function(t,e){var r=e.context;if(t.sol()&&(null==r.align&&(r.align=!1),e.indented=t.indentation(),e.startOfLine=!0,"case"==r.type&&(r.type="}")),t.eatSpace())return null;n=null;var i=(e.tokenize||c)(t,e);return"comment"==i||(null==r.align&&(r.align=!0),"{"==n?s(e,t.column(),"}"):"["==n?s(e,t.column(),"]"):"("==n?s(e,t.column(),")"):"case"==n?r.type="case":("}"==n&&"}"==r.type||n==r.type)&&p(e),e.startOfLine=!1),i},indent:function(e,n){if(e.tokenize!=c&&null!=e.tokenize)return t.Pass;var i=e.context,o=n&&n.charAt(0);if("case"==i.type&&/^(?:case|default)\b/.test(n))return e.context.type="}",i.indented;var a=o==i.type;return i.align?i.column+(a?0:1):i.indented+(a?0:r)},electricChars:"{}):",closeBrackets:"()[]{}''\"\"``",fold:"brace",blockCommentStart:"/*",blockCommentEnd:"*/",lineComment:"//"}})),t.defineMIME("text/x-go","go")}(n(51))}}]);
//# sourceMappingURL=57.3ffcfa97.chunk.js.map