$(function() {
	var mkRow = function() {
		var row = $("<tr/>");

		row.append($("<td><input type=\"text\" name=\"Start\" /></td>"));

		var yn = mkYNCheckbox("RepetitionEnabled", false);
		var cell = $("<td/>");
		cell.append(yn.hidden);
		cell.append(yn.checkbox);
		row.append(cell);

		cell = $("<td/>");
		cell.append($("<input type=\"text\" class=\"quant\" name=\"Count\" /><span> </span>"));
		var unitsel = $("<select name=\"Unit\" size=\"0\" />");
		var units = ["Minute", "Hour", "Day", "Week", "Month", "Year"];
		for(i in units) {
			unitsel.append($("<option value=\""+units[i]+"\">"+units[i]+"(s)</option>"));
		}
		cell.append(unitsel);
		row.append(cell);

		yn = mkYNCheckbox("EndEnabled", false);
		cell = $("<td/>");
		cell.append(yn.hidden);
		cell.append(yn.checkbox);
		row.append(cell);

		row.append($("<td><input type=\"text\" name=\"End\" /></td>"));

		attachFocusHandler($("input, select", row));
		return row;
	};

	var mkYNCheckbox = function(name, b) {
		var hidden = $("<input type=\"hidden\" />");
		hidden.prop("value", b ? "yes" : "no");
		hidden.prop("name", name);

		var checkbox = $("<input type=\"checkbox\" />");
		checkbox.prop("checked", b);
		checkbox.change(function() {
			hidden.prop("value", checkbox.prop("checked") ? "yes" : "no");
		});

		return {"hidden": hidden, "checkbox": checkbox};
	};

	$("select.enabler").each(function(i) {
		var self = $(this);
		var yn = mkYNCheckbox(self.prop("name"), self.val() == "yes");

		self.before(yn.hidden);
		self.before(yn.checkbox);
		self.remove();
	});

	var maxSchedules = $("table.schedules tbody tr").length;
	var checkInsRow = function() {
		if($("table.schedules tbody tr").length >= maxSchedules) {
			return;
		}

		$("table.schedules tbody").append(mkRow());
	};

	var attachFocusHandler = function(q) {
		q.focus(function() {
			var myrow = $(this).parent().parent();
			$("table.schedules tbody tr").not(myrow).each(function(i) {
				checkRemoveRow($(this));
			});

			checkInsRow();
		});
	};

	var checkRemoveRow = function(row) {
		if($("input[name=\"Start\"]", row).val() == "") {
			row.remove();
		}
	};

	var updateSchedulesTab = function() {
		$("table.schedules tbody tr").each(function(i) {
			checkRemoveRow($(this));
		});

		checkInsRow();
	};

	attachFocusHandler($("table.schedules input, table.schedules select"));

	updateSchedulesTab();
});