
$(function () {
	var $tbl = $('#maintable>tbody');

	// Filter - sent each time when list of users is loaded
	var filter = { 
		'search': '',
		'page': 1,
		'deleted': false,
	}

	// load list of users
	function update() {
		$tbl.listLoad({
			url: ajaxUrl + 'getAll',
			method: 'POST',
			noitemsTpl: 'noitems',
			itemTpl: 'item',
			data: filter,
			success: updatePages
		});
	}
	update();

	// Sort
	$.tools.addSorterFn(function (name, direction) {
		filter.sort = (direction === -1 ? '-' : '') + name;
		update();
	});

	// Search
	$.tools.addSearchFn(function (value) {
		filter.search = value;
		update();
	});

	// Ajax tabs
	$('div.tabs').on('change', function () {
		filter.deleted = $('div.tabs').data('active')[0] === "deleted";
		update();
	});

	// "New item" button
	var $newAdmin = $.tools.addButton({
		html: '<span class="fa fa-plus"></span> Add record',
		onClick: function () {
			$.wbox.open('Add new record', window.tplRet('form-add'));
		}
	});

	// Submit form "New item"
	$('#add-form').ajaxFormSender({
		url: ajaxUrl + 'add',
		success: function (result) {
			$tbl.tplPrepend('item', result.data);
			$('#noitems').hide().remove();
			$tbl.children('tr[data-id=' + result.data.id + ']').children().highlight(500);
			return true;
		}
	});

	// "Edit" button pressed
	$.tools.ajaxActionSender('.edit', {
		url: ajaxUrl + 'get',
		method: 'POST',
		success: function (result) {
			$.wbox.open('Change record', window.tplRet('form-edit', result.data));
		}
	});

	// Submit form "Edit"
	$('#edit-form').ajaxFormSender({
		url: ajaxUrl + 'edit',
		success: function (result) {
			result.data.teasersCount = '...';
			$tbl.children('tr[data-id=' + result.data.id + ']').tplReplace('item', result.data);
			$tbl.children('tr[data-id=' + result.data.id + ']').children().highlight(500);
			return true;
		}
	});

	// "Remove/Restore" button pressed
	$.tools.ajaxActionSender('.remove, .restore', {
		url: ajaxUrl + 'action',
		method: 'POST',
		remove: true // remove from list if success
	});

	// Pagination
	function updatePages(data) {
		var pages = data.count / data.limit;
		$('#page-count').text(Math.ceil(pages) || 1);
		$('#page-current').text(Math.ceil(data.page) || 1);
		$('#page-next').attr('disabled', 'disabled');
		$('#page-before').attr('disabled', 'disabled');
		if (data.page < pages) {
			$('#page-next').removeAttr('disabled');
		}
		if (data.page > 1) {
			$('#page-before').removeAttr('disabled');
		}
	}
	$.tools.forceClick('#page-next', function () {
		filter.page++;
		update(filter);
	});
	$.tools.forceClick('#page-before', function () {
		filter.page--;
		update(filter);
	});

});
