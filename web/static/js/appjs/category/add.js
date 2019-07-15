$().ready(function() {
	validateRule();
});

$.validator.setDefaults({
	submitHandler : function() {
		save();
	}
});
function save() {
	$.ajax({
		cache : true,
		type : "POST",
		url : "/category/save",
		data : $('#categoryAddForm').serialize(),// 你的formid
		async : false,
		error : function(request) {
			parent.layer.alert("Connection error");
		},
		success : function(data) {
			if (data.code == 0) {
				parent.layer.msg("操作成功");
				parent.reLoad();
				var index = parent.layer.getFrameIndex(window.name); // 获取窗口索引
				parent.layer.close(index);

            } else if (data.code == -1){
                var msg = "";
                for(var key in data.msg){
                    msg += data.msg[key];
                    msg += "<BR>";
                }
                parent.layer.msg(msg);
            } else {
                parent.layer.msg(data.msg);
            }

		}
	});

}
function validateRule() {
	var icon = "<i class='fa fa-times-circle'></i> ";
	$("#categoryAddForm").validate({
		// rules : {
		// 	name : {
		// 		required : true
		// 	}
		// },
		// messages : {
		// 	name : {
		// 		required : icon + "请输入姓名"
		// 	}
		// }
	})
}