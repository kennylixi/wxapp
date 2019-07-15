// 以下为官方示例
$().ready(function() {
	validateRule();
	// $("#signupForm").validate();
});

$.validator.setDefaults({
	submitHandler : function() {
		update();
	}
});

$(function() {
    Simditor.locale = 'zh-CN';//设置中文
    var editor = new Simditor({
        textarea: $('#Content'),  //textarea的id
        placeholder: '这里输入文字...',
        toolbar: ['title', 'bold', 'italic', 'underline', 'strikethrough', 'fontScale', '|', 'ol', 'ul', 'blockquote', 'table', '|', 'image', 'hr', '|', 'indent', 'outdent', 'alignment'], //工具条都包含哪些内容
        pasteImage: true,//允许粘贴图片
        defaultImage: "http://static.139ud.com/front/img/nav-logo.png",
        upload: {
            url: '/upload/simditor', //文件上传的接口地址
            params: null, //键值对,指定文件上传接口的额外参数,上传的时候随文件一起提交
            fileKey: 'file', //服务器端获取文件数据的参数名
            connectionCount: 3,
            leaveConfirm: '正在上传文件'
        }
    });
});

function update() {
	$.ajax({
		cache : true,
		type : "POST",
		url : "/evaluate/update",
		data : $('#signupForm').serialize(),// 你的formid
		async : false,
		error : function(request) {
			alert("Connection error");
		},
		success : function(data) {
			if (data.code == 0) {
				parent.layer.msg(data.msg);
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