ethereum = require("../ethereum.json")

let mConfig = {}
ethereum.methods.map(m => mConfig[m.name] = {
    Summary: m.summary,
    ResultName: m.result.name,
    Description: m.result.description
})

let jstring = JSON.stringify(mConfig, null, 4).replaceAll(`"Summary"`, "Summary").replaceAll(`"ResultName"`, "ResultName").replaceAll(`"Description"`, "Description")
console.log(jstring)