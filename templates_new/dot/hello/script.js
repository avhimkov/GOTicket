
$(function () {
	$.tools.addLink({
		html: '<i class="fa fa-github"></i> Source Code',
		href: 'https://github.com/night-codes/summer'
	});
	$.tools.addLink({
		html: '<i class="fa fa-info-circle"></i> API Documentation',
		href: 'https://godoc.org/github.com/night-codes/summer'
	});
	$.tools.addButton({
		html: '<i class="fa fa-gift" aria-hidden="true"></i> Icons list',
		id: "button",
		onClick: function () {
			$.wbox.open('Icons used in Summer', window.tplRet('icons'));
		}
	});
	$('.glyph').shwark('icoinfo', {
		data: function (el) {
			return {
				'class': $(el).attr('class')
			};
		}
	});
});
