export const pureImpl = function(a) { return a; };
export const bindImpl = function(a) { return function(b) { return a; }; };
export const liftF = function(a) { return a; };
export const resumePrime = function(a) { return function(b) { return function(c) { return a; }; }; };
export const bindNodeClass = null;
export const bindLeafClass = null;
export const freeObjClass = null;
