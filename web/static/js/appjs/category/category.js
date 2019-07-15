var prefix = "/category"
$(function() {
	load();
});

function load() {
	$('#exampleTable')
		.bootstrapTreeTable(
			{
				id : 'Id',
				code : 'Id',
                parentCode : 'Pid',
				type : "GET", // 请求数据的ajax类型
				url : prefix + '/list', // 请求数据的ajax的url
				ajaxParams : {}, // 请求数据的ajax的data属性
				expandColumn : '1', // 在哪一列上面显示展开按钮
				striped : true, // 是否各行渐变色
				bordered : true, // 是否显示边框
				expandAll : false, // 是否全部展开
				// toolbar : '#exampleToolbar',
				columns : [
					{
						title : '编号',
						field : 'Id',
						visible : false,
						align : 'center',
						valign : 'center',
						width : '60px',
						checkbox : true
					},
					{
						field : 'Title',
						title : '类目名称',
                        align : 'center',
                        valign : 'center',
					},
                    {
                        field : 'Short',
                        title : '类目拼音',
                        align : 'center',
                        valign : 'center',
                    },
                    {
                        field : 'ShortTitle',
                        title : '类目简称',
                        align : 'center',
                        valign : 'center',
                    },
					{
						field : 'Sort',
						title : '排序',
                        align : 'center',
                        valign : 'center',
					},
					{
						field : 'Status',
						title : '状态',
						align : 'center',
                        valign : 'center',
						formatter : function(item, index) {
							if (item.Status == 2) {
								return '<span class="label label-danger">禁用</span>';
							} else if (item.Status == 1) {
								return '<span class="label label-primary">正常</span>';
							}
						}
					},
					{
						title : '操作',
						field : 'id',
						align : 'center',
                        valign : 'center',
						formatter : function(item, index) {
							var e = '<a class="btn btn-primary btn-sm" href="#" mce_href="#" title="编辑" onclick="edit(\''
								+ item.Id
								+ '\')"><i class="fa fa-edit"></i></a> ';
							var a = '<a class="btn btn-primary btn-sm" href="#" title="增加下級"  mce_href="#" onclick="add(\''
								+ item.Id
								+ '\')"><i class="fa fa-plus"></i></a> ';
							var d = '<a class="btn btn-warning btn-sm" href="#" title="删除"  mce_href="#" onclick="remove(\''
								+ item.Id
								+ '\')"><i class="fa fa-remove"></i></a> ';
							return e + a + d;
						}
					} ]
			});
}
function reLoad() {
	load();
}
function add(pId) {
	layer.open({
		type : 2,
		title : '增加',
		maxmin : true,
		shadeClose : false, // 点击遮罩关闭层
		area : [ '800px', '520px' ],
		content : prefix + '/add/' + pId
	});
}
function edit(id) {
	layer.open({
		type : 2,
		title : '编辑',
		maxmin : true,
		shadeClose : false, // 点击遮罩关闭层
		area : [ '800px', '520px' ],
		content : prefix + '/edit/' + id // iframe的url
	});
}
function remove(id) {
	layer.confirm('确定要删除选中的记录？', {
		btn : [ '确定', '取消' ]
	}, function() {
		$.ajax({
			url : prefix + "/remove",
			type : "post",
			data : {
				'id' : id
			},
			success : function(r) {
				if (r.code == 0) {
					layer.msg(r.msg);
					reLoad();
				} else {
					layer.msg(r.msg);
				}
			}
		});
	})
}