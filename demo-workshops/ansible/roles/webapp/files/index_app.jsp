<%@taglib prefix="c" uri="http://java.sun.com/jsp/jstl/core"%>
<%@taglib prefix="sql" uri="http://java.sun.com/jsp/jstl/sql"%>
<%--
    Document   : index
    Created on : Jul 9, 2013, 1:06:39 PM
    Author     : vtigadi


http://127.0.0.1:8080/index_mysql_lc.jsp?dataSource=jdbc%2Fmysql
jdbc/mysql

    <%=session.getAttribute("sql_select_employees")%>

--%>

<%@ include file="init.jsp" %>
<%@ include file="include_page_calcs.jsp" %>

<%
String sql = "";
session.setAttribute("employee_id","NO");
if (session.getAttribute("dbType").equals("oracle")) {
   sql = "select c.* from (";
   sql = sql + "select a.*, count(*) over () as cnt from (";
   sql = sql + "select b.*, rownum as rnum from (";
   sql = sql + "select * from employees" + sorder + " " + sortname + " " + sortdirection;
   sql = sql + ") b";
   sql = sql + ") a";
   sql = sql + ") c";
   sql = sql + " where rnum between " + istart + " and " + iend + "";
   if (isearch != null && !isearch.equals("")) {
      isearch = isearch.toLowerCase();
      sql = "select c.* from (";
      sql = sql + "select a.*, count(*) over () as cnt from (";
      sql = sql + "select b.*, rownum as rnum from (";
      sql = sql + "select * from employees";
      sql = sql + " WHERE lower(first_name) like '%"+isearch+"%'";
      sql = sql + " OR lower(last_name) like '%"+isearch+"%'";
      sql = sql + " OR lower(dept_name) like '%"+isearch+"%'";
      sql = sql + " OR lower(city) like '%"+isearch+"%'";
      sql = sql + sorder + " " + sortname + " " + sortdirection;
      sql = sql + ") b";
      sql = sql + ") a";
      sql = sql + ") c";
      sql = sql + " where rnum between " + istart + " and " + iend + "";
   }
} else {
   sql = (String) session.getAttribute("sql_select_employees");
}
//out.println("SQL> "+sql+"<br />");
%>

<sql:query var="employees" dataSource="<%=session.getAttribute(\"dataSource\")%>">
<%=sql%>
</sql:query>

<%@page contentType="text/html" pageEncoding="UTF-8"%>
<%@page errorPage="ShowError.jsp" %>
<!DOCTYPE html>
<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
        <%--
        <link rel="stylesheet" type="text/css" href="style.css">
        <link href="./bootstrap/css/bootstrap1.min.css" rel="stylesheet" media="screen"/>
        --%>
        <link rel="stylesheet" href="<%=request.getContextPath()%>/bootstrap/css/bootstrap1.min.css" media="screen"/>
<%
String title=pageContext.getServletContext().getInitParameter("myTitle");
%>
        <title><c:out value="<%=title%>" escapeXml="false"/></title>
        <style type="text/css">
            body {
                padding-top: 20px;
                padding-bottom: 20px;
            }

            .alert {
            color: red;
            text-align:center;
            }

            /* Custom container */
            .container {
                margin: 0 auto;
                max-width: 1000px;
            }
            .container > hr {
                margin: 60px 0;
            }

            /* Main marketing message and sign up button */
            .jumbotron {
                margin: 80px 0;
                text-align: center;
            }
            .jumbotron h1 {
                font-size: 100px;
                line-height: 1;
            }
            .jumbotron .lead {
                font-size: 24px;
                line-height: 1.25;
            }
            .jumbotron .btn {
                font-size: 21px;
                padding: 14px 24px;
            }

            /* Supporting marketing content */
            .marketing {
                margin: 60px 0;
            }
            .marketing p + h4 {
                margin-top: 28px;
            }

            /* Customize the navbar links to be fill the entire space of the .navbar */
            .navbar .navbar-inner {
                padding: 0;
            }
            .navbar .nav {
                margin: 0;
                display: table;
                width: 100%;
            }
            .navbar .nav li {
                display: table-cell;
                width: 0%;
                float: none;
            }
            .navbar .nav li a {
                font-weight: bold;
                text-align: center;
                border-left: 1px solid rgba(255,255,255,.75);
                border-right: 1px solid rgba(0,0,0,.1);
            }
            .navbar .nav li:first-child a {
                border-left: 0;
                border-radius: 3px 0 0 3px;
            }
            .navbar .nav li:last-child a {
                border-right: 0;
                border-radius: 0 3px 3px 0;
            }

        </style>
        <link rel="stylesheet" href="<%=request.getContextPath()%>/bootstrap/css/bootstrap-responsive.min.css" media="screen"/>
        <%--<link href="./bootstrap/css/bootstrap-responsive.min.css" rel="stylesheet" media="screen"/> --%>

    </head>
    <body>

    <a name="top"></a>
    <%=session.getAttribute("html_banner")%>

        <script>
            function validateInput()
            {
                var empid = document.getElementById("empid").value;
                var firstname = document.getElementById("firstname").value;
                var lastname = document.getElementById("lastname").value;
                var deptname = document.getElementById("deptname").value;
                var city = document.getElementById("city").value;
                var twitter_handle = document.getElementById("twitter_handle").value;

                if (empid == "" || firstname == "" || lastname == "" || deptname == "" || city == "")
                {
                    bootbox.alert("One or more input fields are empty!");
                    return false;
                }
            }
        </script>
        <h2  class="alert" id="testing"><c:out value="<%=title%>" escapeXml="false"/></h2>
        <div class="container">
            <div class="row">
                <div class="span10">
                    <div class="masthead">
                        <h3 class="muted">Delphix Demo </h3>
                        <div class="navbar" hidden=true>
                            <div class="navbar-inner">
                                <div class="container">
                                    <ul class="nav">
                                        <li class="active"><a href="index_app.jsp?sessionid=<%=sessionid%>">Home</a></li>
                                        <li><a href="masking_app.jsp?sessionid=<%=sessionid%>">Masking</a><li>
                                        <li><a href="login.jsp">Self Service</a></li>
                                    </ul>
                                </div>
                            </div>
                        </div><!-- /.navbar -->
                    </div>
                </div>
                <div class="span2">
                </div>
            </div>
        </div>

        <div class='container'>


<table border=0 width="100%" style="border-collapse:collapse;"><tr><td><h3>Employees</h3></td>
<%@ include file="include_page_html.jsp" %>
<td>&nbsp;&nbsp;&nbsp;&nbsp;</td><!--td align="right"><a href="#top"><image src="images/back2top-icon-2.gif" height="30px" border=0 alt="Back to Top" /></a></td--></tr></table>


            <div class="row">
                <div class="span10">

                    <table id="" class="table table-striped">
                        <!-- column headers -->
                        <tr>
                            <c:forEach var="columnName" items="${employees.columnNames}">
<c:set var="myTest" value="${columnName}"/>
<%
String name=pageContext.getAttribute("myTest").toString();
if (name.equals("EMPLOYEE_ID")) {
   session.setAttribute("employee_id","YES");
}
// out.println("empid "+session.getAttribute("employee_id")+"<br />");
%>
<%@ include file="include_page_sort.jsp"%>
                                </c:forEach>
                        </tr>
                        <!-- column data -->
                        <c:forEach var="row" items="${employees.rows}">
                            <tr>
                                <c:forEach var="columnName" items="${employees.columnNames}">
<c:set var="myTest" value="${columnName}"/>
<c:set var="myVal" value="${row[columnName]}"/>

<%
String strStr=pageContext.getAttribute("myTest").toString();
String strVal = pageContext.getAttribute("myVal").toString();
if (strStr.equals("CNT")) {
   //out.println("yeah"+strVal+":");
   itotal = Integer.parseInt(strVal);
}
if ( !strStr.equals("RNUM") && !strStr.equals("CNT") ) {
               String ztmp = strVal;
               if (isearch != null && !isearch.equals("") ) {
                  if (ztmp != null) {
                     if (ztmp.toLowerCase().contains(isearch)) {
                        ztmp = "<font color=red>"+ztmp+"</font>";
                     }
                  }
               }
//${row[columnName]}
%>
                                   <td><c:out value="<%=ztmp%>" escapeXml="false"/></td>
<%
}
%>
                                </c:forEach>
                            </tr>
                        </c:forEach>

                    </table>
                    <br/>
                    <br/>
                </div>
                <div class="span2">
                </div>
            </div>
        </div>

        <script src="<%=request.getContextPath()%>/bootstrap/js/bootstrap.min.js"></script>
        <script src="<%=request.getContextPath()%>/bootstrap/js/bootbox.min.js"></script>

<%@ include file="footer.jsp" %>