let birthDate = new Date(1997, 0, 3);
let ageDifMs = Date.now() - birthDate.getTime();
let ageDate = new Date(ageDifMs);
let my_age = Math.abs(ageDate.getUTCFullYear() - 1970);
document.getElementById("age").innerHTML = my_age;
