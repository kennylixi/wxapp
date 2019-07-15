$().ready(function() {
	validateRule();
});

$.validator.setDefaults({
	submitHandler : function() {
		save();
	}
});

$('#Province').on('change', function() {
    selectLoad( $('#Province option:selected').val() )
})

function selectLoad(id) {
    var html = "";
    $('#City').empty();
    if(id===""){
        html += '<option value="">-请选择-</option>'
        $('#City').append(html);
        return
    }
    $.ajax({
        url: '/city/city/' + id,
        success: function (data) {
            //加载数据
            html += '<option value="">-请选择-</option>'
            for (var i = 0; i < data.length; i++) {
                html += '<option value="' + data[i].Id + '">' + data[i].Name + '</option>'
            }
            $('#City').append(html);
        }
    });
}

function save() {
	$.ajax({
		cache : true,
		type : "POST",
		url : "/apply/save",
		data : $('#signupForm').serialize(),// 你的formid
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
	$("#signupForm").validate({
		rules : {

		},
		messages : {

		}
	})
}

var openCategory = function(){
	layer.open({
		type:2,
		title:"选择类目",
		area : [ '300px', '450px' ],
		content:"/category/treeview"
	})
}
function loadCategory( id, name ){
	$("#Cid").val(id);
	$("#Cname").val(name);
}