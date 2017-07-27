/**
 * SalesPlace.js
 *
 * @description :: TODO: You might write a short summary of how this model works and what it represents here.
 * @docs        :: http://sailsjs.org/documentation/concepts/models-and-orm/models
 */

module.exports = {

  attributes: {
	name : {
		type : 'string'		
	},
	inventory : {
		type : 'integer'
	}
/* 	stock : {
		collection : 'product'
		via:	'sp_sotck'
	},
	sales_list : {
		collection : 'sales'
		via:	'sp_sale'
	},
	clients_list : {
		collection : 'client'
		via:	'sp_client'
	},
	sellers_list : {
		collection : 'sellers'
		via:	'sp_sellers'
	} */

  }
};

