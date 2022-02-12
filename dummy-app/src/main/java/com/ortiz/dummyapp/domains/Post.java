package com.ortiz.dummyapp.domains;


import java.util.Objects;

public class Post {

  private Long id;
  private Long userId;
  private String title;
  private String body;

  public Long getId() {
    return id;
  }

  public void setId(Long id) {
    this.id = id;
  }

  public Long getUserId() {
    return userId;
  }

  public void setUserId(Long userId) {
    this.userId = userId;
  }

  public String getTitle() {
    return title;
  }

  public void setTitle(String title) {
    this.title = title;
  }

  public String getBody() {
    return body;
  }

  public void setBody(String body) {
    this.body = body;
  }

  @Override
  public boolean equals(Object o) {
    if (this == o) return true;
    if (o == null || getClass() != o.getClass()) return false;
    Post post = (Post) o;
    return Objects.equals(id, post.id) && Objects.equals(userId, post.userId) && Objects.equals(title, post.title) && Objects.equals(body, post.body);
  }

  @Override
  public int hashCode() {
    return Objects.hash(id, userId, title, body);
  }
}
