//const signer = require('nacl-signature');
var nacl=require("tweetnacl")
nacl.util = require('tweetnacl-util');
//var decodeUtf8 = require('decode-utf8')
let keys = nacl.sign.keyPair();
//let publicKey = encoding.toHexString(keys.publicKey).toUpperCase();
//let secretKey = encoding.toHexString(keys.secretKey).toUpperCase();
//this.publicKey.value = publicKey;
//this.secretKey.value = secretKey;
var pub=Buffer.from(keys.publicKey).toString('hex').toUpperCase()
var priv=Buffer.from(keys.secretKey).toString('hex').toUpperCase()
var k=nacl.util.decodeUTF8("kui")
var signature=nacl.sign(k, keys.secretKey)
console.log("Public Key: " + Buffer.from(keys.publicKey).toString('hex').toUpperCase() + "\n")
console.log("Private Key: " + Buffer.from(keys.secretKey).toString('hex').toUpperCase() + "\n")
var sign = Buffer.from(signature).toString('hex').toUpperCase()
console.log(sign.substr(0,128))
console.log("Signature: " + Buffer.from(signature).toString('hex').toUpperCase() + "\n")

//console.log(signature)
//console.log())
//console.log(signature)
