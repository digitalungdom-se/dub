module.exports = {
  name: 'factorial',
  description: 'Räknar ut fakultet',
  aliases: ['fakultet','fk','fc'],
  group: 'misc',
  usage: 'factorial <command>',
  execute(message, args) {

  	if (!args.length || isNaN(args[0]) || args[0] < 0) {
  		return message.reply('that\'s not a valid command!');
  	}
  	
  	const factorial = n => n ? n * factorial(--n) : 1;

  	try{return message.reply(factorial(args[0]))}catch{return message.reply('The number is too big')};
}
};
