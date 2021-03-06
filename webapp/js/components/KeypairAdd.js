/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */
'use strict'

var React = require('react');
var injectIntl = require('react-intl').injectIntl;
var Keypairs = require('../models/keypairs');
var Navigation = require('./Navigation');
var AlertBox = require('./AlertBox');
var Footer = require('./Footer');

var KeypairAdd = React.createClass({
	getInitialState: function() {
    return {authorityId: null, key: null, error: this.props.error};
  },

	handleChangeAuthorityId: function(e) {
		this.setState({authorityId: e.target.value});
	},

	handleChangeKey: function(e) {
		this.setState({key: e.target.value});
	},

	handleFileUpload: function(e) {
		var self = this;
    var reader = new FileReader();
    var file = e.target.files[0];

    reader.onload = function(upload) {
      self.setState({
        key: upload.target.result.split(',')[1],
      });
    }

    reader.readAsDataURL(file);
	},

	handleSaveClick: function(e) {
		var self = this;
		e.preventDefault();

		Keypairs.create(this.state.authorityId, this.state.key).then(function(response) {
			var data = JSON.parse(response.body);
			if ((response.statusCode >= 300) || (!data.success)) {
        self.setState({error: self.formatError(data)});
      } else {
        window.location = '/models';
      }
		});
	},

	formatError: function(data) {
		var message = this.props.intl.formatMessage({id: data.error_code});
		if (data.error_subcode) {
			message += ': ' + this.props.intl.formatMessage({id: data.error_subcode});
		} else if (data.message) {
			message += ': ' + data.message;
		}
		return message;
	},

	render: function() {
		var M = this.props.intl.formatMessage;

		return (
			<div>
				<Navigation active="models" />

				<section className="row no-border">
					<h2>{M({id: 'new-signing-key'})}</h2>
					<div className="twelve-col">
						<AlertBox message={this.state.error} />

						<form>
							<fieldset>
								<ul>
									<li>
										<label htmlFor="authority-id">{M({id: 'authority-id'})}:</label>
										<input type="text" id="authority-id" onChange={this.handleChangeAuthorityId} placeholder={M({id: 'authority-id-description'})} />
									</li>
									<li>
										<label htmlFor="key">{M({id: 'signing-key'})}:</label>
										<textarea onChange={this.handleChangeKey} defaultValue={this.state.key} id="key"
												placeholder={M({id: 'new-signing-key-description'})}>
										</textarea>
									</li>
									<li>
										<input type="file" onChange={this.handleFileUpload} />
									</li>
								</ul>
							</fieldset>
						</form>
						<div>
							<a href='/models' onClick={this.handleSaveClick} className="button--primary">{M({id: 'save'})}</a>
							&nbsp;
							<a href='/models' className="button--secondary">{M({id: 'cancel'})}</a>
						</div>
					</div>
				</section>

				<Footer />
			</div>
		);
	}
});

module.exports = injectIntl(KeypairAdd);
