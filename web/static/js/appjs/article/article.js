var prefix = "/article"
var categoryId = '';
$(function() {
	getTreeData([1,2]);
	load();
});

function load() {
	$('#exampleTable')
		.bootstrapTable(
			{
                cache : false,
				method : 'get', // 服务器数据的请求方式 get or post
				url : prefix + "/list", // 服务器数据的加载地址
				// showRefresh : true,
				// showToggle : true,
				// showColumns : true,
				iconSize : 'outline',
				toolbar : '#exampleToolbar',
				striped : true, // 设置为true会有隔行变色效果
				dataType : "json", // 服务器返回的数据类型
				pagination : true, // 设置为true会在底部显示分页条
				// queryParamsType : "limit",
				// //设置为limit则会发送符合RESTFull格式的参数
				singleSelect : false, // 设置为true将禁止多选
				// contentType : "application/x-www-form-urlencoded",
				// //发送到服务器的数据编码类型
				pageSize : 10, // 如果设置了分页，每页数据条数
				pageNumber : 1, // 如果设置了分布，首页页码
				// search : true, // 是否显示搜索框
				showColumns : false, // 是否显示内容下拉框（选择显示的列）
				sidePagination : "server", // 设置在哪里进行分页，可选值为"client" 或者
				// "server"
				queryParams : function(params) {
					return {
						// 说明：传入后台的参数包括offset开始索引，limit步长，sort排序列，order：desc或者,以及所有列的键值对
						limit : params.limit,
						offset : params.offset,
						Title : $('#searchName').val(),
                        Cid : categoryId,
					};
				},
				// //请求服务器数据时，你可以通过重写参数的方式添加一些额外的参数，例如 toolbar 中的参数 如果
				// queryParamsType = 'limit' ,返回参数必须包含
				// limit, offset, search, sort, order 否则, 需要包含:
				// pageSize, pageNumber, searchText, sortName,
				// sortOrder.
				// 返回false将会终止请求
				columns : [
					{
						checkbox : true
					},
					{
						field : 'Id', // 列字段名
						title : '序号' // 列标题
					},
					{
						field : 'Title',
						title : '标题'
					},
					{
                        field : 'Status',
                        title : '状态',
                        align : 'center',
                        formatter : function(value, row, index) {
                            if (value == '2') {
                                return '<span class="label label-danger">已下架</span>';
                            } else if (value == '1') {
                                return '<span class="label label-primary">已上架</span>';
                            }
                        }
					},
					{
						title : '操作',
						field : 'id',
						align : 'center',
						formatter : function(value, row, index) {
							var d = '<a class="btn btn-warning btn-sm" href="#" title="删除"  mce_href="#" onclick="remove(\''
								+ row.Id
								+ '\')"><i class="fa fa-remove"></i></a> ';
                            var f = '<a class="btn btn-success btn-sm" href="#" title="上下架"  mce_href="#" onclick="grounding(\''
                                + row.Id
                                + '\')"><i class="fa fa-key"></i></a> ';
                            var g = '<a class="btn btn-success btn-sm" href="#" title="加入专题"  mce_href="#" onclick="addSpecial(\''
                                + row.Id
                                + '\')"><i class="fa fa-plus"></i></a> ';
							return d + f + g;
						}
					} ]
			});
}
function reLoad() {
	$('#exampleTable').bootstrapTable('refresh');
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
function grounding(id) {
    layer.confirm('确定要上下架选中的记录？', {
        btn : [ '确定', '取消' ]
    }, function() {
        $.ajax({
            url : prefix + "/grounding",
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
function addSpecial(id) {
    layer.open({
        type : 2,
        title : '加入专题',
        maxmin : true,
        shadeClose : false,
        area : [ '800px', '520px' ],
        content : prefix + '/special/' + id // iframe的url
    });
}
function batchRemove() {
	var rows = $('#exampleTable').bootstrapTable('getSelections'); // 返回所有选择的行，当没有选择的记录时，返回一个空数组
	if (rows.length == 0) {
		layer.msg("请选择要删除的数据");
		return;
	}
	layer.confirm("确认要删除选中的'" + rows.length + "'条数据吗?", {
		btn : [ '确定', '取消' ]
	// 按钮
	}, function() {
		var ids = new Array();
		// 遍历所有选择的行数据，取每条数据对应的ID
		$.each(rows, function(i, row) {
			ids[i] = row['Id'];
		});
		$.ajax({
			type : 'POST',
			data : {
				"ids" : JSON.stringify(ids)
			},
			url : prefix + '/removes',
			success : function(r) {
				if (r.code == 0) {
					layer.msg(r.msg);
					reLoad();
				} else {
					layer.msg(r.msg);
				}
			}
		});
	}, function() {});
}
function getTreeData(ids) {
	$.ajax({
		type : "GET",
		url : "/category/tree",
        data : {
            "ids" : JSON.stringify(ids)
        },
		success : function(tree) {
			loadTree(tree);
		}
	});
}
function loadTree(tree) {
	$('#jstree').jstree({
		'core' : {
			'data' : tree
		},
		"plugins" : [ "search" ]
	});
	$('#jstree').jstree().open_all();
}
$('#jstree').on("changed.jstree", function(e, data) {
	if (data.selected == -1) {
        categoryId = ''
	} else {
        categoryId = data.selected[0]
	}
    var opt = {
        query : {
            Cid : categoryId,
            Title : $('#searchName').val(),
            Status : $('#Status option:selected').val(),
        }
    }
    $('#exampleTable').bootstrapTable('refresh',opt);
});

$('#Status').on('change', function() {
    var opt = {
        query : {
            Cid : categoryId,
            Title : $('#searchName').val(),
            Status : $('#Status option:selected').val(),
        }
    }
    $('#exampleTable').bootstrapTable('refresh', opt);
});