module.exports = {
  name: 'vote',
  description: 'Rösta om en fråga',
  aliases: ['vote','vt','rösta'],
  group: 'misc',
  usage: 'rösta <id> <alternativindex>',
  example: '',
  serverOnly: true,
  execute(message, args) {
  	if ( args.length === 0 ) return message.reply( 'du måste ge ett id' );
  	if ( args.length === 1 ) return message.reply( 'du måste ge ett index' );
  	const id = args[0];
  	const index = args[1];
  	const authorId = message.author.id;
  	if ( !voteDic.hasOwnProperty(id)){return message.reply( 'det finns ingen pågående röstning med detta id' ); }

  	if ( isNaN(index) || index < 0 || index >= voteDic[ id ][ 'options' ].length) return message.reply( 'du måste ge ett existerande index' );

  	if ( voteDic[ id ][ 'users' ][ authorId ] !== undefined ) voteDic[ id ][ 'score' ][ voteDic[ id ][ 'users' ][ authorId ] ]--;
  	if ( voteDic[ id ][ 'score' ][ index ] === undefined ) voteDic[ id ][ 'score' ][ index ] = 0;


	message.delete();
	voteDic[id][ 'users' ][ authorId ] = index;
	voteDic[id][ 'score' ][ index ]++;
	message.reply( voteDic[ id ] )
	console.log( voteDic );
  
  },

};
