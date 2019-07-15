var prefix = "/apply"
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
                        field : 'Pass',
                        title : '状态',
                        align : 'center',
                        formatter : function(value, row, index) {
                            if (value == '3') {
                                return '<span class="label label-primary">已通过</span>';
                            } else if (value == '2') {
                                return '<span class="label label-danger">未通过</span>';
                            } else if (value == '1') {
                                return '<span class="label label-warn">未审核</span>';
                            }
                        }
					},
					{
						title : '操作',
						field : 'id',
						align : 'center',
						formatter : function(value, row, index) {
							var e = '<a  class="btn btn-primary btn-sm" href="#" mce_href="#" title="编辑" onclick="edit(\''
								+ row.Id
								+ '\')"><i class="fa fa-edit "></i></a> ';
                            var d = '<a class="btn btn-success btn-sm" href="#" title="审核"  mce_href="#" onclick="pass(\''
                                + row.Id
                                + '\')"><i class="fa fa-key"></i></a> ';
							return e + d;
						}
					} ]
			});
}
function reLoad() {
	$('#exampleTable').bootstrapTable('refresh');
}
function add() {
	// iframe层
	layer.open({
		type : 2,
		title : '增加文章',
		maxmin : true,
		shadeClose : false, // 点击遮罩关闭层
		area : [ '800px', '520px' ],
		content : prefix + '/add'
	});
}
function pass(id) {
    layer.open({
        type : 2,
        title : '审核',
        maxmin : true,
        shadeClose : false, // 点击遮罩关闭层
        area : [ '400px', '260px' ],
        content : prefix + '/pass/' + id // iframe的url
    });
}
function edit(id) {
	layer.open({
		type : 2,
		title : '修改文章',
		maxmin : true,
		shadeClose : false,
		area : [ '800px', '520px' ],
		content : prefix + '/edit/' + id // iframe的url
	});
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
            pass : $('#Pass option:selected').val(),
        }
    }
    $('#exampleTable').bootstrapTable('refresh',opt);
});

$('#Pass').on('change', function() {
    var opt = {
        query : {
            Cid : categoryId,
            Title : $('#searchName').val(),
            pass : $('#Pass option:selected').val(),
        }
    }
    $('#exampleTable').bootstrapTable('refresh', opt);
});